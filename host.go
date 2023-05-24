package ep2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	maddr "github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

type Host struct {
	ctx context.Context
	dht *DHT
	h   host.Host

	bootstrapPeers []maddr.Multiaddr
	logger         *zap.Logger
}

func NewHost(bootstrapPeers []maddr.Multiaddr, logger *zap.Logger) *Host {
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	host := &Host{
		ctx:            context.Background(),
		h:              h,
		bootstrapPeers: bootstrapPeers,
		logger:         logger,
	}
	host.dht = newDHT(host, logger)

	return host
}

func (h *Host) ID() peer.ID {
	return h.h.ID()
}

func (h *Host) Desc() string {
	return fmt.Sprintf("%v/p2p/%v", h.h.Addrs(), h.h.ID())
}

func (h *Host) Peers() []peer.ID {
	return h.h.Network().Peers()
}

func (h *Host) NewGossipSub() *GossipSub {
	return newGossipSub(h, h.logger)
}
