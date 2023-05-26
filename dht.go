package ep2p

import (
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	// dual "github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
	"go.uber.org/zap"
)

type DHT struct {
	// kademliaDHT      *dual.DHT
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery

	host   *Host
	logger *zap.Logger
}

func newDHT(host *Host, logger *zap.Logger, serverMode bool) *DHT {
	logger = logger.Named("DHT")
	// For test.
	// var kademliaDHT *dual.DHT
	// var err error
	// if !serverMode {
	// 	kademliaDHT, err = dual.New(host.ctx, host.h, dual.DHTOption(dht.Mode(dht.ModeClient)))
	// } else {
	// 	kademliaDHT, err = dual.New(host.ctx, host.h, dual.DHTOption(dht.Mode(dht.ModeServer)))
	// }
	var kademliaDHT *dht.IpfsDHT
	var err error
	if !serverMode {
		kademliaDHT, err = dht.New(host.ctx, host.impl, dht.Mode(dht.ModeClient), dht.ProtocolPrefix("/ep2p/dht/1.0"))
	} else {
		kademliaDHT, err = dht.New(host.ctx, host.impl, dht.Mode(dht.ModeServer), dht.ProtocolPrefix("/ep2p/dht/1.0"))
	}
	// kademliaDHT, err := dht.New(host.ctx, host.h, dht.Mode(dht.ModeAutoServer))
	if err != nil {
		panic(err)
	}

	if err = kademliaDHT.Bootstrap(host.ctx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range host.bootstrapPeers {
		wg.Add(1)
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		go func() {
			defer wg.Done()
			if err := host.impl.Connect(host.ctx, *peerInfo); err != nil {
				logger.Info("failed to connect to bootstrap peer", zap.Error(err))
			} else {
				logger.Info("connected to bootstrap peer", zap.String("peer", peerInfo.String()))
			}
		}()
	}
	wg.Wait()

	return &DHT{
		kademliaDHT:      kademliaDHT,
		routingDiscovery: routing.NewRoutingDiscovery(kademliaDHT),
		host:             host,
		logger:           logger,
	}
}

func (dht *DHT) Advertise(ns string) error {
	util.Advertise(dht.host.ctx, dht.routingDiscovery, ns)
	return nil
}

func (dht *DHT) FindPeers(rendezvous string) (<-chan peer.AddrInfo, error) {
	return dht.routingDiscovery.FindPeers(dht.host.ctx, rendezvous)
}
