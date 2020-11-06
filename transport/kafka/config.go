package kafka

import (
	"emperror.dev/errors"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/events"
)

type Config struct {
	events.Config `mapstructure:",squash"`
	Addr          []string
	GroupName     string
	SubOverwrite  *sarama.Config
	PubOverwrite  *sarama.Config
	TopicDetails  *sarama.TopicDetail
}

func (c Config) SetDefaults(env config.Env) {
	*c.PubOverwrite = *(kafka.DefaultSaramaSyncPublisherConfig())
	*c.SubOverwrite = *(kafka.DefaultSaramaSubscriberConfig())
}

func (c Config) Validate() (err error) {
	if len(c.Addr) == 0 {
		err = errors.Append(err, errors.New("kafka: addr is empty"))
	}

	return
}
