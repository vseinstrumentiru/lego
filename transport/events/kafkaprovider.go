package events

import (
	"github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"
	"go.uber.org/dig"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/v2/multilog"
	kafkatransport "github.com/vseinstrumentiru/lego/v2/transport/kafka"
)

type KafkaPubArgs struct {
	dig.In
	Broker  *kafkatransport.Config
	Logger  multilog.Logger
	Encoder kafka.Marshaler
}

type KafkaSubArgs struct {
	dig.In
	Broker  *kafkatransport.Config
	Logger  multilog.Logger
	Decoder kafka.Unmarshaler
}

func ProvideKafkaPublisher(in KafkaPubArgs) (*kafka.Publisher, error) {
	pub, err := kafka.NewPublisher(
		in.Broker.Addr,
		in.Encoder,
		in.Broker.PubOverwrite,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "kafka", "component": "events.publisher"})),
	)

	return pub.(*kafka.Publisher), err
}

func ProvideKafkaSubscriber(in KafkaSubArgs) (*kafka.Subscriber, error) {
	subCfg := kafka.SubscriberConfig{
		Brokers:                in.Broker.Addr,
		ConsumerGroup:          in.Broker.GroupName,
		NackResendSleep:        *in.Broker.AckTimeout,
		InitializeTopicDetails: in.Broker.TopicDetails,
	}

	if in.Broker.AckTimeout != nil && *in.Broker.AckTimeout > 0 {
		subCfg.NackResendSleep = *in.Broker.AckTimeout
	}

	if in.Broker.ReconnectTimeout != nil && *in.Broker.ReconnectTimeout > 0 {
		subCfg.ReconnectRetrySleep = *in.Broker.ReconnectTimeout
	}

	sub, err := kafka.NewSubscriber(
		subCfg,
		in.Broker.SubOverwrite,
		in.Decoder,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"provider": "kafka", "component": "events.subscriber"})),
	)

	return sub.(*kafka.Subscriber), err
}
