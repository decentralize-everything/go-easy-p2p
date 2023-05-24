package main

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"strings"

	ep2p "github.com/decentralize-everything/go-easy-p2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
	topicName      = "test-topic"
)

func main() {
	flag.Var(&bootstrapPeers, "peers", "bootstrap peer list")
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	host := ep2p.NewHost(bootstrapPeers, logger)
	logger.Info("host created", zap.String("desc", host.Desc()))

	gs := host.NewGossipSub()
	topic, err := gs.Join(
		topicName,
		ep2p.WithCallback(func(m *pubsub.Message) error {
			from := peer.ID(m.Message.From).String()
			receivedFrom := m.ReceivedFrom.String()
			logger.Info(
				"new message received",
				zap.String("origin", from[len(from)-6:]),
				zap.String("sender", receivedFrom[len(receivedFrom)-6:]),
				zap.String("topic", *m.Message.Topic),
				zap.String("data", string(m.Message.Data)),
			)
			return errors.New("error for test")
		}),
		ep2p.FilterSelf(host.ID()),
	)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if s == "quit\n" {
			gs.Leave(topicName)
			return
		}

		if err := topic.Publish([]byte(s)); err != nil {
			panic(err)
		}

		logger.Info("peers connected", zap.Any("peers", host.Peers()))
	}
}
