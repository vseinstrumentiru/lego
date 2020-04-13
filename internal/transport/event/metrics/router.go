package metrics

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"time"
)

func Register(r *message.Router) {
	r.AddPublisherDecorators(DecoratePublisher)
	r.AddSubscriberDecorators(DecorateSubscriber)
	r.AddMiddleware(Middleware)
}

func Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		now := time.Now()
		ctx := msg.Context()

		tags := []tag.Mutator{
			tag.Upsert(HandlerName, message.HandlerNameFromCtx(ctx)),
		}

		defer func() {
			if err != nil {
				tags = append(tags, tag.Upsert(Success, "false"))
			} else {
				tags = append(tags, tag.Upsert(Success, "true"))
			}

			_ = stats.RecordWithTags(ctx, tags, HandlerExecutionTime.M(float64(time.Since(now))/float64(time.Millisecond)))
		}()

		return h(msg)
	}
}
