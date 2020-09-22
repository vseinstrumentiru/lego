package metrics

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"time"
)

var (
	publisherLabelKeys = []string{
		labelKeyHandlerName,
		labelKeyPublisherName,
		labelSuccess,
	}
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
	labels := labelsFromCtx(ctx, publisherLabelKeys...)
	if labels[labelKeyPublisherName] == "" {
		labels[labelKeyPublisherName] = p.publisherName
	}
	if labels[labelKeyHandlerName] == "" {
		labels[labelKeyHandlerName] = labelValueNoHandler
	}

	tags := []tag.Mutator{
		tag.Upsert(PublisherName, labels[labelKeyPublisherName]),
		tag.Upsert(HandlerName, labels[labelKeyHandlerName]),
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

		_ = stats.RecordWithTags(ctx, tags, PublisherPublishTime.M(time.Since(start).Seconds()))
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
func DecoratePublisher(pub message.Publisher) (message.Publisher, error) {
	return &PublisherDecorator{
		pub:           pub,
		publisherName: StructName(pub),
	}, nil
}
