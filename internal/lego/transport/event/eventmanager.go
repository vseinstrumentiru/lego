package event

import (
	"context"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	originNats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sagikazarmark/kitx/correlation"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/transport/event/metrics"
	"github.com/vseinstrumentiru/lego/tools/eventtools/cloudevent"
	watermilllog "logur.dev/integration/watermill"
	"os"
	"regexp"
	"sync"
	"time"
)

type publishers struct {
	sync.Mutex
	defaultName string
	defaultPub  message.Publisher
	items       map[string]message.Publisher
}

func (em *publishers) add(key, name string, publisher message.Publisher) (err error) {
	var pubErr error
	publisher, pubErr = metrics.DecoratePublisher(name, publisher)
	err = errors.Append(err, pubErr)
	publisher, pubErr = message.MessageTransformPublisherDecorator(func(msg *message.Message) {
		if cid, ok := correlation.FromContext(msg.Context()); ok {
			middleware.SetCorrelationID(cid, msg)
		}
	})(publisher)

	if key == em.defaultName {
		em.defaultPub = publisher
	}

	em.items[key] = publisher

	return
}

func (em *publishers) Publish(topic string, messages ...*message.Message) error {
	var pub message.Publisher
	var ok bool

	em.Lock()
	if pub, ok = em.items[topic]; !ok {
		pub = em.defaultPub
	}
	em.Unlock()

	if pub == nil {
		return errors.New("undefined publisher")
	}

	return pub.Publish(topic, messages...)
}

func (em *publishers) Close() (err error) {
	for _, pub := range em.items {
		err = errors.Append(err, pub.Close())
	}

	return
}

type subscribers struct {
	sync.Mutex
	defaultName string
	defaultSub  message.Subscriber
	items       map[string]message.Subscriber
}

func (em *subscribers) add(key, name string, subscriber message.Subscriber) (err error) {
	var subErr error
	subscriber, subErr = metrics.DecorateSubscriber(name, subscriber)
	err = errors.Append(err, subErr)
	subscriber, subErr = message.MessageTransformSubscriberDecorator(func(msg *message.Message) {
		if cid := middleware.MessageCorrelationID(msg); cid != "" {
			msg.SetContext(correlation.ToContext(msg.Context(), cid))
		}
	})(subscriber)

	if key == em.defaultName {
		em.defaultSub = subscriber
	}

	em.items[key] = subscriber

	return
}

func (em *subscribers) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	var sub message.Subscriber
	var ok bool

	em.Lock()
	if sub, ok = em.items[topic]; !ok {
		sub = em.defaultSub
	}
	em.Unlock()

	return sub.Subscribe(ctx, topic)
}

func (em *subscribers) Close() (err error) {
	for _, pub := range em.items {
		err = errors.Append(err, pub.Close())
	}

	return
}

type eventManager struct {
	lego.LogErr
	*publishers
	*subscribers
	router *message.Router
}

func newEventManager(logErr lego.LogErr, config Config) (_ *eventManager, err error) {
	em := &eventManager{
		LogErr: logErr,
		publishers: &publishers{
			defaultName: config.DefaultProvider,
			items:       make(map[string]message.Publisher),
		},
		subscribers: &subscribers{
			defaultName: config.DefaultProvider,
			items:       make(map[string]message.Subscriber),
		},
	}

	for name, cfg := range config.Providers.Nats {
		var marshaller nats.MarshalerUnmarshaler
		if cfg.CloudEvent.Enabled {
			marshaller = cloudevent.Marshaller{
				SpecVersion: "0.3",
				Source:      cfg.CloudEvent.Source,
			}
		} else {
			marshaller = nats.GobMarshaler{}
		}

		var suffixer func(string) string

		switch cfg.ClientIDSuffixGen {
		case natsClientIDSuffixHost:
			suffixer = hostSuffix
		default:
			suffixer = withoutSuffix
		}

		if cfg.Pub {
			var pub message.Publisher
			var pubErr error
			natsConn, conErr := originNats.Connect(cfg.Addr, natsConnOptions(logErr, cfg.PanicOnLost)...)

			if err != nil {
				err = errors.Append(err, conErr)
			}
			pub, pubErr = nats.NewStreamingPublisher(
				nats.StreamingPublisherConfig{
					ClusterID: cfg.ClusterID,
					ClientID:  suffixer(cfg.ClientID + "_pub"),
					StanOptions: []stan.Option{
						stan.NatsConn(natsConn),
						stan.Pings(20, 10),
					},
					Marshaler: marshaller,
				},
				watermilllog.New(logErr.WithFields(map[string]interface{}{"component": "events.nats.pub." + name})),
			)

			err = errors.Append(err, pubErr)
			err = errors.Append(err, em.publishers.add(name, "events.nats.pub."+name, pub))
		}

		if cfg.Sub {
			natsConn, conErr := originNats.Connect(cfg.Addr, natsConnOptions(logErr, cfg.PanicOnLost)...)

			if err != nil {
				err = errors.Append(err, conErr)
			}

			sub, subErr := nats.NewStreamingSubscriber(
				nats.StreamingSubscriberConfig{
					ClusterID:        cfg.ClusterID,
					ClientID:         suffixer(cfg.ClientID + "_sub"),
					QueueGroup:       cfg.QueueGroup,
					DurableName:      cfg.DurableName,
					SubscribersCount: cfg.SubscribersCount, // how many goroutines should consume messages
					CloseTimeout:     cfg.CloseTimeout,
					AckWaitTimeout:   cfg.AckWaitTimeout,
					StanOptions: []stan.Option{
						stan.NatsConn(natsConn),
						stan.Pings(20, 10),
					},
					StanSubscriptionOptions: []stan.SubscriptionOption{
						stan.DeliverAllAvailable(),
					},
					Unmarshaler: marshaller,
				},
				watermilllog.New(logErr.WithFields(map[string]interface{}{"component": "events.nats.sub." + name})),
			)

			err = errors.Append(err, subErr)

			err = errors.Append(err, subErr)
			err = errors.Append(err, em.subscribers.add(name, "events.nats.sub."+name, sub))
		}
	}

	for name, cfg := range config.Providers.Channel {
		pubsub := gochannel.NewGoChannel(
			gochannel.Config{
				OutputChannelBuffer:            cfg.OutputChannelBuffer,
				Persistent:                     cfg.Persistent,
				BlockPublishUntilSubscriberAck: cfg.BlockPublishUntilSubscriberAck,
			},
			watermilllog.New(logErr.WithFields(map[string]interface{}{"component": "events.channel.pubsub." + name})),
		)
		err = errors.Append(err, em.publishers.add(name, "events.channel.pub."+name, pubsub))
		err = errors.Append(err, em.subscribers.add(name, "events.channel.sub."+name, pubsub))
	}

	{
		router, routerErr := message.NewRouter(
			message.RouterConfig{
				CloseTimeout: config.RouterConfig.CloseTimeout,
			},
			watermilllog.New(logErr.WithFields(map[string]interface{}{"component": "events.router"})),
		)

		err = errors.Append(err, routerErr)

		retryMiddleware := middleware.Retry{}
		retryMiddleware.MaxRetries = config.RouterConfig.MaxRetries
		retryMiddleware.MaxInterval = config.RouterConfig.MaxRetryInterval

		router.AddMiddleware(
			// if retries limit was exceeded, message is sent to poison queue (poison_queue topic)
			retryMiddleware.Middleware,

			// correlation ID middleware adds to every produced message correlation id of consumed message,
			// useful for debugging
			middleware.CorrelationID,
		)

		em.router = router
	}

	return em, err
}

func (e *eventManager) AddHandlers(
	handlers []cqrs.EventHandler,
	generateTopic func(eventName string) string,
	marshaler cqrs.CommandEventMarshaler,
) error {
	processor, err := cqrs.NewEventProcessor(
		handlers,
		generateTopic,
		func(handlerName string) (message.Subscriber, error) { return e.subscribers, nil },
		marshaler,
		watermilllog.New(e.LogErr.WithFields(map[string]interface{}{"component": "events.processor"})),
	)

	if err != nil {
		return err
	}

	return processor.AddHandlersToRouter(e.router)
}

func (e *eventManager) Publisher() message.Publisher {
	return e.publishers
}

func (e *eventManager) Run(ctx context.Context) (err error) {
	return e.router.Run(ctx)
}

func (e *eventManager) Close() (err error) {
	err = errors.Append(err, e.subscribers.Close())
	err = errors.Append(err, e.publishers.Close())
	return
}

func (e *eventManager) Subscriber() message.Subscriber {
	return e.subscribers
}

func withoutSuffix(str string) string {
	return str
}

func hostSuffix(str string) string {
	name, err := os.Hostname()
	if err != nil {
		emperror.Panic(err)
	}

	return str + "_" + regexp.MustCompile(`[^a-zA-Z0-9_\-]`).ReplaceAllString(name, "_")
}

func natsConnOptions(log lego.LogErr, panicOnLostConn bool, opts ...originNats.Option) []originNats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, originNats.ReconnectWait(reconnectDelay))
	opts = append(opts, originNats.MaxReconnects(int(totalWait/reconnectDelay)))
	if panicOnLostConn {
		opts = append(opts, originNats.DisconnectErrHandler(func(nc *originNats.Conn, err error) {
			emperror.Panic(errors.Wrap(err, "nats disconnected"))
		}))
		opts = append(opts, originNats.ClosedHandler(func(nc *originNats.Conn) {
			log.Error("Exiting, no servers available")
		}))
	} else {
		opts = append(opts, originNats.DisconnectErrHandler(func(nc *originNats.Conn, err error) {
			log.Info(fmt.Sprintf("Disconnected: will attempt reconnects for %.0fm", totalWait.Minutes()))
		}))
		opts = append(opts, originNats.ReconnectHandler(func(nc *originNats.Conn) {
			log.Info(fmt.Sprintf("Reconnected [%s]", nc.ConnectedUrl()))
		}))
		opts = append(opts, originNats.ClosedHandler(func(nc *originNats.Conn) {
			log.Error("Exiting, no servers available")
		}))
	}

	return opts
}
