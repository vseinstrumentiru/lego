package events

import (
	"github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/events"
	"github.com/vseinstrumentiru/lego/multilog"
	kafkatransport "github.com/vseinstrumentiru/lego/transport/kafka"
)

type kafkaArgs struct {
	Kafka  *kafkatransport.Config
	Logger multilog.Logger
	Config *events.Config `optional:"true"`
}

func ProvideKafka(in kafkaArgs) (events.KafkaProvider, error) {
	p := kafkaProvider{
		brokers: in.Kafka.Addr,
		logger:  in.Logger.WithFields(map[string]interface{}{"provider": "kafka"}),
	}

	if in.Config != nil {
		p.config = *in.Config
	}

	return p, nil
}

type kafkaProvider struct {
	brokers []string
	config  events.Config
	logger  multilog.Logger
}

func (k kafkaProvider) NewPublisher(cfg events.KafkaPublisherConfig) (message.Publisher, error) {
	return kafka.NewPublisher(
		k.brokers,
		cfg.Encoder,
		cfg.Config,
		watermilllog.New(k.logger.WithFields(map[string]interface{}{"component": "events.publisher"})),
	)
}

func (k kafkaProvider) NewSubscriber(cfg events.KafkaSubscriberConfig) (message.Subscriber, error) {
	subCfg := kafka.SubscriberConfig{
		Brokers:                k.brokers,
		ConsumerGroup:          cfg.GroupName,
		NackResendSleep:        cfg.AckTimeout,
		InitializeTopicDetails: cfg.TopicDetails,
	}

	if subCfg.ReconnectRetrySleep <= 0 && k.config.ReconnectTimeout != nil {
		subCfg.ReconnectRetrySleep = *k.config.ReconnectTimeout
	}

	return kafka.NewSubscriber(
		subCfg,
		cfg.Overwrite,
		cfg.Decoder,
		watermilllog.New(k.logger.WithFields(map[string]interface{}{"component": "events.subscriber"})),
	)
}
