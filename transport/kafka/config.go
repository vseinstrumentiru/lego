package kafka

import (
	"github.com/Shopify/sarama"

	"github.com/vseinstrumentiru/lego/events"
)

type Config struct {
	events.Config
	Addr         []string
	GroupName    string
	Overwrite    *sarama.Config
	TopicDetails *sarama.TopicDetail
}
