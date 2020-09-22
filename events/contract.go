package events

import (
	"context"

	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Router interface {
	AddMiddleware(m ...message.HandlerMiddleware)
	AddPlugin(p ...message.RouterPlugin)
	AddPublisherDecorators(dec ...message.PublisherDecorator)
	AddSubscriberDecorators(dec ...message.SubscriberDecorator)
	AddHandler(
		handlerName string,
		subscribeTopic string,
		subscriber message.Subscriber,
		publishTopic string,
		publisher message.Publisher,
		handlerFunc message.HandlerFunc,
	) *message.Handler
	AddNoPublisherHandler(
		handlerName string,
		subscribeTopic string,
		subscriber message.Subscriber,
		handlerFunc message.NoPublishHandlerFunc,
	)
	Run(ctx context.Context) (err error)
	Running() chan struct{}
	Close() error
}

type TransportEncoder interface {
	Encode(topic string, msg *message.Message) (payload []byte, err error)
}

type TransportDecoder interface {
	Decode(payload []byte, headers map[string]string) (msg *message.Message, err error)
}

type MessageEncoder interface {
	Encode(v interface{}) (*message.Message, error)
}

type MessageDecoder interface {
	Decode(msg *message.Message, v interface{}) (err error)
}

type StanProvider interface {
	NewPublisher(m nats.Marshaler) (message.Publisher, error)
	NewSubscriber(cfg NatsSubscriberConfig) (message.Subscriber, error)
}

type KafkaProvider interface {
	NewPublisher(cfg KafkaPublisherConfig) (message.Publisher, error)
	NewSubscriber(cfg KafkaSubscriberConfig) (message.Subscriber, error)
}

type ChannelProvider interface {
	NewPublisher() (message.Publisher, error)
	NewSubscriber() (message.Subscriber, error)
}
