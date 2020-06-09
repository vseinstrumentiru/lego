package metrics

import (
	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

type SubscriberDecorator struct {
	message.Subscriber
	sub            message.Subscriber
	subscriberName string
}

func (s *SubscriberDecorator) recordMetrics(msg *message.Message) {
	if msg == nil {
		return
	}

	ctx := msg.Context()

	subscriberName := message.SubscriberNameFromCtx(ctx)
	if subscriberName == "" {
		subscriberName = s.subscriberName
	}

	handlerName := message.HandlerNameFromCtx(ctx)
	if handlerName == "" {
		handlerName = tagValueNoHandler
	}

	tags := []tag.Mutator{
		tag.Upsert(SubscriberName, subscriberName),
		tag.Upsert(HandlerName, handlerName),
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

func (s *SubscriberDecorator) Close() error {
	return s.sub.Close()
}

// DecorateSubscriber decorates a publisher with instrumentation.
func DecorateSubscriber(name string, sub message.Subscriber) (message.Subscriber, error) {
	d := &SubscriberDecorator{
		sub:            sub,
		subscriberName: name,
	}

	var err error
	d.Subscriber, err = message.MessageTransformSubscriberDecorator(d.recordMetrics)(sub)
	if err != nil {
		return nil, errors.Wrap(err, "could not decorate subscriber with metrics decorator")
	}

	return d, nil
}
