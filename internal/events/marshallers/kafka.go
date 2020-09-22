package marshallers

import (
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/vseinstrumentiru/lego/events"
)

type KafkaEncoder struct {
	events.TransportEncoder
}

func (k KafkaEncoder) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	panic("implement me")
}
