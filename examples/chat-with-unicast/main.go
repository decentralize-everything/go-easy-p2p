package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	ep2p "github.com/decentralize-everything/go-easy-p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	maddr "github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

type addrList []maddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

var (
	bootstrapPeers addrList
)

func main() {
	flag.Var(&bootstrapPeers, "boots", "bootstrap peer list")
	target := flag.String("dst", "", "target peer to send")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var host *ep2p.Host
	// if *target != "" {
	if len(bootstrapPeers) != 0 {
		host = ep2p.NewClient(bootstrapPeers, logger)
	} else {
		host = ep2p.NewServer(bootstrapPeers, logger)
	}
	// host := ep2p.NewServer(bootstrapPeers, logger)

	logger.Info("host created", zap.Strings("desc", host.Desc()))
	if len(bootstrapPeers) == 0 {
		logger.Info("run as bootstrap node")
		select {}
	}

	if *target == "" {
		host.SetDefaultRecvHandler(func(s network.Stream) {
			buf := bufio.NewReader(s)
			for {
				str, err := buf.ReadString('\n')
				if err != nil {
					panic(err)
				}

				logger.Info("message received", zap.String("message", string(str)))
			}
		})
		logger.Info("run as receiver node")
		select {}
	}

	logger.Info("run as sender node")
	receiver, err := peer.Decode(*target)
	if err != nil {
		panic(err)
	}

	logger.Info("peers connected before send", zap.Any("peers", host.Peers()))
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if s == "quit\n" {
			return
		}

		if err := host.Send([]byte(s), receiver); err != nil {
			panic(err)
		}

		logger.Info("peers connected", zap.Any("peers", host.Peers()))
	}
}
