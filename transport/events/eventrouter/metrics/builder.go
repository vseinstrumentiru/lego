package metrics

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

func Register(r *message.Router) {
	r.AddPublisherDecorators(DecoratePublisher)
	r.AddSubscriberDecorators(DecorateSubscriber)
	r.AddMiddleware(Middleware)
}
