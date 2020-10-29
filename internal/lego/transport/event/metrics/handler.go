package metrics

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

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

			_ = stats.RecordWithTags(ctx, tags, HandlerExecutionTime.M(time.Since(now).Seconds()))
		}()

		return h(msg)
	}
}
