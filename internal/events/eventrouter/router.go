package eventrouter

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/multilog"
)

func newRouter(cfg message.RouterConfig, logger multilog.Logger) (*router, error) {
	r, err := message.NewRouter(cfg, watermilllog.New(logger))
	if err != nil {
		return nil, err
	}

	return &router{
		hasHandlers: false,
		i:           r,
	}, nil
}

type router struct {
	hasHandlers bool
	i           *message.Router
}

func (r *router) AddMiddleware(m ...message.HandlerMiddleware) {
	r.i.AddMiddleware(m...)
}

func (r *router) AddPlugin(p ...message.RouterPlugin) {
	r.i.AddPlugin(p...)
}

func (r *router) AddPublisherDecorators(dec ...message.PublisherDecorator) {
	r.i.AddPublisherDecorators(dec...)
}

func (r *router) AddSubscriberDecorators(dec ...message.SubscriberDecorator) {
	r.i.AddSubscriberDecorators(dec...)
}

func (r *router) AddHandler(handlerName string, subscribeTopic string, subscriber message.Subscriber, publishTopic string, publisher message.Publisher, handlerFunc message.HandlerFunc) *message.Handler {
	return r.i.AddHandler(
		handlerName,
		subscribeTopic,
		subscriber,
		publishTopic,
		publisher,
		handlerFunc,
	)
}

func (r *router) AddNoPublisherHandler(handlerName string, subscribeTopic string, subscriber message.Subscriber, handlerFunc message.NoPublishHandlerFunc) {
	r.i.AddNoPublisherHandler(handlerName, subscribeTopic, subscriber, handlerFunc)
}

func (r *router) Run(ctx context.Context) (err error) {
	if !r.hasHandlers {
		return nil
	}

	return r.i.Run(ctx)
}

func (r *router) Running() chan struct{} {
	return r.i.Running()
}

func (r *router) Close() error {
	return r.i.Close()
}
