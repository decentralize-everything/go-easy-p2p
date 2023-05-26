package ep2p

import (
	"context"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.uber.org/zap"
)

type GossipSubOption func(*Topic) error

func WithCallback(callback GossipSubCallback) GossipSubOption {
	return func(t *Topic) error {
		t.callback = callback
		return nil
	}
}

func FilterSelf(self peer.ID) GossipSubOption {
	return func(t *Topic) error {
		t.filterSelf = true
		t.self = self
		return nil
	}
}

type GossipSubCallback func(m *pubsub.Message) error

type Topic struct {
	ctx   context.Context
	name  string
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	callback   GossipSubCallback
	filterSelf bool
	self       peer.ID

	logger *zap.Logger
}

func newTopic(ctx context.Context, name string, topic *pubsub.Topic, sub *pubsub.Subscription, logger *zap.Logger, opts ...GossipSubOption) *Topic {
	t := &Topic{
		ctx:    ctx,
		name:   name,
		topic:  topic,
		sub:    sub,
		logger: logger.With(zap.String("name", name)),
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.callback != nil {
		go t.subscribeRoutine()
	}

	return t
}

func (t *Topic) subscribeRoutine() {
	for {
		m, err := t.sub.Next(t.ctx)
		if err != nil {
			panic(err)
		}

		if t.filterSelf && peer.ID(m.Message.From) == t.self {
			continue
		}

		if err = t.callback(m); err != nil {
			t.logger.Warn("failed to handle message in callback", zap.Any("message", m), zap.Error(err))
			continue
		}
	}
}

func (t *Topic) Publish(data []byte) error {
	if err := t.topic.Publish(t.ctx, data); err != nil {
		return Libp2pError(err)
	}
	return nil
}

type GossipSub struct {
	topics map[string]*Topic
	ps     *pubsub.PubSub
	host   *Host
	logger *zap.Logger
}

func newGossipSub(host *Host, logger *zap.Logger) *GossipSub {
	ps, err := pubsub.NewGossipSub(host.ctx, host.impl)
	if err != nil {
		panic(err)
	}

	return &GossipSub{
		topics: make(map[string]*Topic),
		ps:     ps,
		host:   host,
		logger: logger.Named("GossipSub"),
	}
}

func (gs *GossipSub) Join(topic string, opts ...GossipSubOption) (*Topic, error) {
	if _, found := gs.topics[topic]; found {
		return nil, TopicAlreadyExistsError(topic)
	}

	gs.host.dht.Advertise(topic)
	go gs.discoverPeers(topic)
	t, err := gs.ps.Join(topic)
	if err != nil {
		return nil, Libp2pError(err)
	}
	sub, err := t.Subscribe()
	if err != nil {
		return nil, Libp2pError(err)
	}

	return newTopic(gs.host.ctx, topic, t, sub, gs.logger, opts...), nil
}

func (gs *GossipSub) Leave(topic string) error {
	if t, found := gs.topics[topic]; !found {
		return TopicNotFoundError(topic)
	} else {
		t.sub.Cancel()
		return nil
	}
}

func (gs *GossipSub) discoverPeers(rendezvous string) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-gs.host.ctx.Done():
			return
		case <-ticker.C:
			peers, err := gs.host.dht.FindPeers(rendezvous)
			if err != nil {
				panic(err)
			}

			for peer := range peers {
				if peer.ID == gs.host.impl.ID() {
					continue
				}

				if gs.host.impl.Network().Connectedness(peer.ID) == network.Connected {
					continue
				}

				if err := gs.host.impl.Connect(gs.host.ctx, peer); err != nil {
					gs.logger.Info("failed to connect to peer", zap.String("peer", peer.ID.Pretty()), zap.Error(err))
				} else {
					gs.logger.Info("connected to peer", zap.String("peer", peer.ID.Pretty()))
				}
			}
		}
	}
}
