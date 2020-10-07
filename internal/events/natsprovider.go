package events

import (
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/nats-io/stan.go"
	"go.uber.org/dig"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/multilog"
	lestan "github.com/vseinstrumentiru/lego/transport/stan"
)

type natsPubArgs struct {
	dig.In
	Stan    stan.Conn
	Logger  multilog.Logger
	Encoder nats.Marshaler
}

type natsSubArgs struct {
	dig.In
	Config  lestan.Config
	Stan    stan.Conn
	Logger  multilog.Logger
	Decoder nats.Unmarshaler
}

func ProvideNatsPublisher(in natsPubArgs) (*nats.StreamingPublisher, error) {
	return nats.NewStreamingPublisherWithStanConn(
		in.Stan,
		nats.StreamingPublisherPublishConfig{
			Marshaler: in.Encoder,
		},
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "stan", "component": "events.publisher"})),
	)
}

func ProvideNatsSubscriber(in natsSubArgs) (*nats.StreamingSubscriber, error) {
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

	return nats.NewStreamingSubscriberWithStanConn(
		in.Stan,
		subCfg,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "stan", "component": "events.subscriber"})),
	)
}
