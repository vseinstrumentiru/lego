package events

import (
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/events"
	"github.com/vseinstrumentiru/lego/multilog"
)

type natsArgs struct {
	Config *events.Config `optional:"true"`
	Stan   stan.Conn
	Logger multilog.Logger
}

type natsProvider struct {
	events.Config
	conn   stan.Conn
	logger multilog.Logger
}

func ProvideNats(in natsArgs) events.StanProvider {
	c := natsProvider{
		conn:   in.Stan,
		logger: in.Logger.WithFields(map[string]interface{}{"provider": "stan"}),
	}

	if in.Config != nil {
		c.Config = *in.Config
	}

	return c
}

func (n natsProvider) NewPublisher(m nats.Marshaler) (message.Publisher, error) {
	return nats.NewStreamingPublisherWithStanConn(
		n.conn,
		nats.StreamingPublisherPublishConfig{
			Marshaler: m,
		},
		watermilllog.New(n.logger.WithFields(map[string]interface{}{"component": "events.publisher"})),
	)
}

func (n natsProvider) NewSubscriber(cfg events.NatsSubscriberConfig) (message.Subscriber, error) {
	subCfg := nats.StreamingSubscriberSubscriptionConfig{
		AckWaitTimeout: cfg.AckTimeout,
		CloseTimeout:   cfg.CloseTimeout,
		DurableName:    cfg.DurableName,
		QueueGroup:     cfg.GroupName,
		Unmarshaler:    cfg.Decoder,
	}

	if subCfg.AckWaitTimeout <= 0 && n.AckWaitTimeout != nil {
		subCfg.AckWaitTimeout = *n.AckWaitTimeout
	}

	if subCfg.CloseTimeout <= 0 && n.CloseTimeout != nil {
		subCfg.CloseTimeout = *n.CloseTimeout
	}

	return nats.NewStreamingSubscriberWithStanConn(
		n.conn,
		subCfg,
		watermilllog.New(n.logger.WithFields(map[string]interface{}{"component": "events.subscriber"})),
	)
}
