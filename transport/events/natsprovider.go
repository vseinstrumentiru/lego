package events

import (
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/nats-io/stan.go"
	"go.uber.org/dig"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/v2/log"
	lestan "github.com/vseinstrumentiru/lego/v2/transport/stan"
)

type NatsPubArgs struct {
	dig.In
	Stan    stan.Conn
	Logger  log.Logger
	Encoder nats.Marshaler
}

type NatsSubArgs struct {
	dig.In
	Config  *lestan.Config
	Stan    stan.Conn
	Logger  log.Logger
	Decoder nats.Unmarshaler
}

func ProvideNatsPublisher(in NatsPubArgs) (*nats.StreamingPublisher, error) {
	pub, err := nats.NewStreamingPublisherWithStanConn(
		in.Stan,
		nats.StreamingPublisherPublishConfig{
			Marshaler: in.Encoder,
		},
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "stan", "component": "events.publisher"})),
	)

	return pub, err
}

func ProvideNatsSubscriber(in NatsSubArgs) (*nats.StreamingSubscriber, error) {
	subCfg := nats.StreamingSubscriberSubscriptionConfig{
		DurableName: in.Config.DurableName,
		QueueGroup:  in.Config.GroupName,
		Unmarshaler: in.Decoder,
	}

	if in.Config.AckTimeout != nil && *in.Config.AckTimeout > 0 {
		subCfg.AckWaitTimeout = *in.Config.AckTimeout
	}

	if in.Config.CloseTimeout != nil && *in.Config.CloseTimeout > 0 {
		subCfg.CloseTimeout = *in.Config.CloseTimeout
	}

	sub, err := nats.NewStreamingSubscriberWithStanConn(
		in.Stan,
		subCfg,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "stan", "component": "events.subscriber"})),
	)

	return sub, err
}
