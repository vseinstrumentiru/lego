package kafka

import (
	"emperror.dev/errors"
	"github.com/Shopify/sarama"

	"github.com/vseinstrumentiru/lego/events"
)

type Config struct {
	events.Config `mapstructure:",squash"`
	Addr          []string
	GroupName     string
	Overwrite     *sarama.Config
	TopicDetails  *sarama.TopicDetail
}

func (c Config) Validate() (err error) {
	if len(c.Addr) == 0 {
		err = errors.Append(err, errors.New("kafka: addr is empty"))
	}

	return
}
