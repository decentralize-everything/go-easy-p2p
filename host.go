package ep2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	// "github.com/libp2p/go-libp2p/core/routing"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	maddr "github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

const (
	defaultStreamProtocol = "/ep2p/unicast/default"
)

type Host struct {
	ctx  context.Context
	dht  *DHT
	impl host.Host

	bootstrapPeers []maddr.Multiaddr
	logger         *zap.Logger

	defaultStreams map[peer.ID]network.Stream
}

func NewServer(bootstrapPeers []maddr.Multiaddr, logger *zap.Logger) *Host {
	return NewHost(bootstrapPeers, logger, true)
}

func NewClient(bootstrapPeers []maddr.Multiaddr, logger *zap.Logger) *Host {
	return NewHost(bootstrapPeers, logger, false)
}

func NewHost(bootstrapPeers []maddr.Multiaddr, logger *zap.Logger, serverMode bool) *Host {
	ho := &Host{
		ctx:            context.Background(),
		bootstrapPeers: bootstrapPeers,
		logger:         logger,
		defaultStreams: make(map[peer.ID]network.Stream),
	}

	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		// libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		// 	ho.h = h
		// 	dht := newDHT(ho, logger, serverMode)
		// 	ho.dht = dht
		// 	return dht.kademliaDHT, nil
		// }),
	)
	if err != nil {
		panic(err)
	}

	ho.impl = h
	ho.dht = newDHT(ho, logger, serverMode)
	ho.impl = routedhost.Wrap(h, ho.dht.kademliaDHT)
	return ho
}

// For testing purposes
func (h *Host) Impl() host.Host {
	return h.impl
}

func (h *Host) ID() peer.ID {
	return h.impl.ID()
}

func (h *Host) Desc() []string {
	var strs []string
	for _, addr := range h.impl.Addrs() {
		strs = append(strs, fmt.Sprintf("%v/p2p/%v", addr, h.impl.ID()))
	}
	return strs
}

func (h *Host) Peers() []peer.ID {
	return h.impl.Network().Peers()
}

func (h *Host) NewGossipSub() *GossipSub {
	return newGossipSub(h, h.logger)
}

func (h *Host) Send(data []byte, peer peer.ID) error {
	if _, found := h.defaultStreams[peer]; !found {
		s, err := h.impl.NewStream(h.ctx, peer, defaultStreamProtocol)
		if err != nil {
			return Libp2pError(err)
		}

		h.defaultStreams[peer] = s
	}

	if _, err := h.defaultStreams[peer].Write(data); err != nil {
		delete(h.defaultStreams, peer)
		return Libp2pError(err)
	}

	return nil
}

func (h *Host) SetDefaultRecvHandler(callback func(network.Stream)) {
	h.impl.SetStreamHandler(defaultStreamProtocol, callback)
}
