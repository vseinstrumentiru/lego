package event

import (
	"context"
	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/sagikazarmark/kitx/correlation"
	"github.com/vseinstrumentiru/lego/internal/transport/event/metrics"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"go.opencensus.io/stats/view"
)

type Server struct {
	Subscriber message.Subscriber
	Publisher  message.Publisher
	Router     *message.Router
}

func (s Server) Close() error {
	defer s.Subscriber.Close()
	defer s.Publisher.Close()
	return nil
}

func Run(p lego.Process, config Config) Server {
	server := Server{}

	var err error
	server.Publisher, server.Subscriber, err = NewPubSub(config, p.Log())
	emperror.Panic(err)

	server.Publisher, _ = metrics.DecoratePublisher(server.Publisher)
	server.Publisher, _ = message.MessageTransformPublisherDecorator(func(msg *message.Message) {
		if cid, ok := correlation.FromContext(msg.Context()); ok {
			middleware.SetCorrelationID(cid, msg)
		}
	})(server.Publisher)

	server.Subscriber, _ = metrics.DecorateSubscriber(server.Subscriber)
	server.Subscriber, _ = message.MessageTransformSubscriberDecorator(func(msg *message.Message) {
		if cid := middleware.MessageCorrelationID(msg); cid != "" {
			msg.SetContext(correlation.ToContext(msg.Context(), cid))
		}
	})(server.Subscriber)

	server.Router, err = NewRouter(config.Router, p.Log())
	emperror.Panic(err)

	metrics.Register(server.Router)
	_ = view.Register(
		metrics.HandlerExecutionTimeView,
		metrics.PublisherPublishTimeView,
		metrics.SubscriberReceivedMessageView,
	)

	p.Run(
		func() error {
			return server.Router.Run(context.Background())
		},
		func(e error) {
			p.Handle(e)
			_ = server.Router.Close()
		})

	return server
}
