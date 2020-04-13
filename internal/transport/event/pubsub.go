package event

import (
	"fmt"
	"github.com/vseinstrumentiru/lego/pkg/eventtools/cloudevent"

	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/nats-io/stan.go"
	watermilllog "logur.dev/integration/watermill"
	"logur.dev/logur"
)

// NewPubSub returns a new PubSub.
func NewPubSub(config Config, logger logur.Logger) (message.Publisher, message.Subscriber, error) {
	switch config.Provider {
	case "nats":
		return NewNatsPubSub(config.Nats, logger)
	case "channel":
		return NewChannelPubSub(config.Channel, logger)
	default:
		return nil, nil, fmt.Errorf("provider %v not allowed", config.Provider)
	}
}

func NewChannelPubSub(config GoChannelProviderConfig, logger logur.Logger) (message.Publisher, message.Subscriber, error) {
	pubsub := gochannel.NewGoChannel(
		gochannel.Config{
			OutputChannelBuffer:            config.OutputChannelBuffer,
			Persistent:                     config.Persistent,
			BlockPublishUntilSubscriberAck: config.BlockPublishUntilSubscriberAck,
		},
		watermilllog.New(logur.WithField(logger, "component", "events.channel")),
	)

	return pubsub, pubsub, nil
}

func NewNatsPubSub(config NatsProviderConfig, logger logur.Logger) (message.Publisher, message.Subscriber, error) {
	var marshaller nats.MarshalerUnmarshaler
	if config.CloudEvent.Enabled {
		marshaller = cloudevent.Marshaller{
			SpecVersion: "0.3",
			Source:      config.CloudEvent.Source,
		}
	} else {
		marshaller = nats.GobMarshaler{}
	}

	subscriber, err := nats.NewStreamingSubscriber(
		nats.StreamingSubscriberConfig{
			ClusterID:        config.ClusterID,
			ClientID:         config.ClientID,
			QueueGroup:       config.QueueGroup,
			DurableName:      config.DurableName,
			SubscribersCount: config.SubscribersCount, // how many goroutines should consume messages
			CloseTimeout:     config.CloseTimeout,
			AckWaitTimeout:   config.AckWaitTimeout,
			StanOptions: []stan.Option{
				stan.NatsURL(config.Addr),
			},
			StanSubscriptionOptions: []stan.SubscriptionOption{
				stan.DeliverAllAvailable(),
			},
			Unmarshaler: marshaller,
		},
		watermilllog.New(logur.WithField(logger, "component", "events.nats")),
	)

	if err != nil {
		return nil, nil, err
	}

	publisher, err := nats.NewStreamingPublisher(
		nats.StreamingPublisherConfig{
			ClusterID: config.ClusterID,
			ClientID:  config.ClientID + "_pub",
			StanOptions: []stan.Option{
				stan.NatsURL(config.Addr),
			},
			Marshaler: marshaller,
		},
		watermilllog.New(logur.WithField(logger, "component", "events.nats")),
	)

	if err != nil {
		return nil, nil, err
	}

	return publisher, subscriber, nil
}
