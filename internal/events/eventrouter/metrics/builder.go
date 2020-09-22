package metrics

import "github.com/vseinstrumentiru/lego/events"

func Register(r events.Router) {
	r.AddPublisherDecorators(DecoratePublisher)
	r.AddSubscriberDecorators(DecorateSubscriber)
	r.AddMiddleware(Middleware)
}
