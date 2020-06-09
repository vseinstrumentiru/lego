package lego

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

type EventManager interface {
	Publisher() message.Publisher
	Subscriber() message.Subscriber
	AddHandlers(handlers []cqrs.EventHandler, generateTopic func(eventName string) string, marshaler cqrs.CommandEventMarshaler) error
	Close() error
}
