package event

import (
	"context"
	"emperror.dev/emperror"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/transport/event/metrics"
	"go.opencensus.io/stats/view"
)

func Run(p lego.Process, config Config) (_ *eventManager, exec func() error, interrupt func(error)) {
	em, err := newEventManager(p, config)
	emperror.Panic(err)

	metrics.Register(em.router)
	_ = view.Register(
		metrics.HandlerExecutionTimeView,
		metrics.PublisherPublishTimeView,
		metrics.SubscriberReceivedMessageView,
	)

	return em,
		func() error { return em.Run(context.Background()) },
		func(e error) {
			p.Handle(e)
			_ = em.router.Close()
		}
}
