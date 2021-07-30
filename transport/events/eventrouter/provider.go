package eventrouter

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/oklog/run"
	"go.opencensus.io/stats/view"
	"go.uber.org/dig"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/v2/events"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/transport/events/eventrouter/metrics"
)

type Args struct {
	dig.In
	RouterConfig *events.RouterConfig `optional:"true"`
	Logger       multilog.Logger
	Pipeline     *run.Group
}

func Provide(in Args) (*message.Router, error) {
	cfg := message.RouterConfig{}
	logger := in.Logger.WithFields(map[string]interface{}{"component": "events.router"})

	var mw []message.HandlerMiddleware
	mw = append(mw, middleware.Recoverer)

	if in.RouterConfig != nil {
		cfg.CloseTimeout = in.RouterConfig.CloseTimeout

		if in.RouterConfig.Retry != nil {
			retryMiddleware := middleware.Retry{
				MaxRetries:  in.RouterConfig.Retry.Count,
				MaxInterval: in.RouterConfig.Retry.Interval,
				Logger:      watermilllog.New(logger),
			}

			mw = append(mw, retryMiddleware.Middleware)
		}

		if in.RouterConfig.Correlation {
			mw = append(mw, middleware.CorrelationID)
		}

		if in.RouterConfig.Throttle != nil {
			mw = append(mw, middleware.NewThrottle(in.RouterConfig.Throttle.Count, in.RouterConfig.Throttle.Interval).Middleware)
		}

		if in.RouterConfig.Timeout != nil {
			mw = append(mw, middleware.Timeout(in.RouterConfig.Timeout.Interval))
		}
	}

	router, err := message.NewRouter(cfg, watermilllog.New(logger))
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(mw...)

	metrics.Register(router)

	err = view.Register(
		metrics.HandlerExecutionTimeView,
		metrics.PublisherPublishTimeView,
		metrics.SubscriberReceivedMessageView,
	)

	in.Pipeline.Add(
		func() error {
			logger.Info("starting router")

			return router.Run(context.Background())
		},
		func(err error) {
			logger.Info("shutting router down")
			_ = router.Close()
		},
	)

	return router, err
}
