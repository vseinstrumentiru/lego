package events

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
)

type Config struct {
	AckWaitTimeout   *time.Duration
	ReconnectTimeout *time.Duration
	CloseTimeout     *time.Duration
}

type NatsSubscriberConfig struct {
	DurableName  string
	GroupName    string
	Count        int
	AckTimeout   time.Duration
	CloseTimeout time.Duration
	Decoder      nats.Unmarshaler
}

type KafkaSubscriberConfig struct {
	GroupName    string
	AckTimeout   time.Duration
	TopicDetails *sarama.TopicDetail
	Decoder      kafka.Unmarshaler
	Overwrite    *sarama.Config
}

type KafkaPublisherConfig struct {
	Encoder kafka.Marshaler
	Config  *sarama.Config
}
