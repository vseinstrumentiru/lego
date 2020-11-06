package deprecated

import (
	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/dig"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

// Deprecated: use LeGo V2
type EventManager interface {
	Publisher() message.Publisher
	Subscriber() message.Subscriber
	AddHandlers(handlers []cqrs.EventHandler, generateTopic func(eventName string) string, marshaler cqrs.CommandEventMarshaler) error
	Close() error
}

type emArgs struct {
	dig.In
	NatsSub   *nats.StreamingSubscriber `optional:"true"`
	NasPub    *nats.StreamingPublisher  `optional:"true"`
	GoChannel *gochannel.GoChannel      `optional:"true"`
	Router    *message.Router
	Log       multilog.Logger
}

func NewEventManager(args emArgs) (EventManager, error) {
	if args.GoChannel == nil && args.NasPub == nil && args.NatsSub == nil {
		return nil, errors.New("event manager: provider not found")
	}

	em := &eventManager{
		router: args.Router,
		log:    args.Log.WithFields(map[string]interface{}{"component": "event_manager"}),
	}

	if args.GoChannel != nil {
		em.pub, em.sub = args.GoChannel, args.GoChannel
	}

	if args.NatsSub != nil {
		em.sub = args.NatsSub
	}

	if args.NasPub != nil {
		em.pub = args.NasPub
	}

	return em, nil
}

type eventManager struct {
	pub    message.Publisher
	sub    message.Subscriber
	router *message.Router
	log    multilog.Logger
}

func (e *eventManager) Publisher() message.Publisher {
	return e.pub
}

func (e *eventManager) Subscriber() message.Subscriber {
	return e.sub
}

func (e *eventManager) AddHandlers(handlers []cqrs.EventHandler, generateTopic func(eventName string) string, marshaler cqrs.CommandEventMarshaler) error {
	processor, err := cqrs.NewEventProcessor(
		handlers,
		generateTopic,
		func(handlerName string) (message.Subscriber, error) { return e.sub, nil },
		marshaler,
		watermilllog.New(e.log.WithFields(map[string]interface{}{"component": "events.processor"})),
	)

	if err != nil {
		return err
	}

	return processor.AddHandlersToRouter(e.router)
}

func (e *eventManager) Close() error {
	return nil
}
