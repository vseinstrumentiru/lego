package metrics

import (
	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

//nolint:gochecknoglobals
var subscriberLabelKeys = []string{
	labelKeyHandlerName,
	labelKeySubscriberName,
}

type SubscriberDecorator struct {
	message.Subscriber
	subscriberName string
}

func (s *SubscriberDecorator) recordMetrics(msg *message.Message) {
	if msg == nil {
		return
	}

	ctx := msg.Context()
	labels := labelsFromCtx(ctx, subscriberLabelKeys...)

	if labels[labelKeySubscriberName] == "" {
		labels[labelKeySubscriberName] = s.subscriberName
	}
	if labels[labelKeyHandlerName] == "" {
		labels[labelKeyHandlerName] = labelValueNoHandler
	}

	tags := []tag.Mutator{
		tag.Upsert(SubscriberName, labels[labelKeySubscriberName]),
		tag.Upsert(HandlerName, labels[labelKeyHandlerName]),
	}

	go func() {
		if subscribeAlreadyObserved(ctx) {
			// decorator idempotency when applied decorator multiple times
			return
		}

		select {
		case <-msg.Acked():
			tags = append(tags, tag.Upsert(Acked, "acked"))
		case <-msg.Nacked():
			tags = append(tags, tag.Upsert(Acked, "nacked"))
		}

		_ = stats.RecordWithTags(ctx, tags, SubscriberReceivedMessage.M(1))
	}()

	msg.SetContext(setSubscribeObservedToCtx(msg.Context()))
}

// DecorateSubscriber decorates a publisher with instrumentation.
func DecorateSubscriber(sub message.Subscriber) (message.Subscriber, error) {
	d := &SubscriberDecorator{
		subscriberName: StructName(sub),
	}

	var err error
	d.Subscriber, err = message.MessageTransformSubscriberDecorator(d.recordMetrics)(sub)
	if err != nil {
		return nil, errors.Wrap(err, "could not decorate subscriber with metrics decorator")
	}

	return d, nil
}
