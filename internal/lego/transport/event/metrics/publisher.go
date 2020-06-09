package metrics

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"time"
)

type PublisherDecorator struct {
	pub           message.Publisher
	publisherName string
}

func (p *PublisherDecorator) Publish(topic string, messages ...*message.Message) (err error) {
	if len(messages) == 0 {
		return p.pub.Publish(topic)
	}

	// TODO: take ctx not only from first msg. Might require changing the signature of Publish, which is planned anyway.
	ctx := messages[0].Context()

	publisherName := message.PublisherNameFromCtx(ctx)
	if publisherName == "" {
		publisherName = p.publisherName
	}

	handlerName := message.HandlerNameFromCtx(ctx)
	if handlerName == "" {
		handlerName = tagValueNoHandler
	}

	tags := []tag.Mutator{
		tag.Upsert(PublisherName, publisherName),
		tag.Upsert(HandlerName, handlerName),
	}

	start := time.Now()

	defer func() {
		if publishAlreadyObserved(ctx) {
			// decorator idempotency when applied decorator multiple times
			return
		}

		if err != nil {
			tags = append(tags, tag.Upsert(Success, "false"))
		} else {
			tags = append(tags, tag.Upsert(Success, "true"))
		}

		_ = stats.RecordWithTags(ctx, tags, PublisherPublishTime.M(float64(time.Since(start))/float64(time.Millisecond)))
	}()

	for _, msg := range messages {
		msg.SetContext(setPublishObservedToCtx(msg.Context()))
	}

	return p.pub.Publish(topic, messages...)
}

func (p *PublisherDecorator) Close() error {
	return p.pub.Close()
}

// DecoratePublisher decorates a publisher with instrumentation.
func DecoratePublisher(name string, pub message.Publisher) (message.Publisher, error) {
	return &PublisherDecorator{
		pub:           pub,
		publisherName: name,
	}, nil
}
