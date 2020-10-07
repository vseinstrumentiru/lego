package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"emperror.dev/emperror"
	"emperror.dev/errors/match"
	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"

	baseCfg "github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/internal/config"
	environment "github.com/vseinstrumentiru/lego/internal/config/env"
	di "github.com/vseinstrumentiru/lego/internal/container"
	"github.com/vseinstrumentiru/lego/internal/events"
	"github.com/vseinstrumentiru/lego/internal/events/eventrouter"
	"github.com/vseinstrumentiru/lego/internal/metrics"
	"github.com/vseinstrumentiru/lego/internal/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/internal/metrics/exporters/prometheus"
	"github.com/vseinstrumentiru/lego/internal/metrics/propagation"
	multilogProvider "github.com/vseinstrumentiru/lego/internal/multilog"
	grpcProvider "github.com/vseinstrumentiru/lego/internal/transpoort/grpc"
	"github.com/vseinstrumentiru/lego/internal/transpoort/http/httpserver"
	"github.com/vseinstrumentiru/lego/multilog"
	"github.com/vseinstrumentiru/lego/server/shutdown"
	"github.com/vseinstrumentiru/lego/version"
)

const (
	defaultConfigPath = "app"
)

func Run(app interface{}, cfg interface{}) {
	// core container instance
	container := newContainer()

	env := environment.New(defaultConfigPath)
	env.SetFlag("version", false, "show version")
	var ver version.Info

	container.
		register(func() di.Container { return container.i }).
		// configuration
		register(func() config.Config { return cfg }).
		register(func() environment.Env { return env }).
		execute(config.Configure).
		execute(func(cfg *baseCfg.Application) {
			ver = version.New(cfg)
		}).
		instance(ver)

	if env.OnFlag("version", func(bool) { ver.Print() }) {
		os.Exit(0)
		return
	}
	// logger
	var logger multilog.Logger
	container.register(multilogProvider.Provide).
		execute(func(l multilog.Logger) {
			logger = l
		})

	defer emperror.HandleRecover(logger.WithFields(map[string]interface{}{"type": "panic"}))
	multilog.SetStandardLogger(logger.WithFields(map[string]interface{}{"type": "standard"}))

	// pipeline configuration
	closer := new(shutdown.CloseGroup)
	container.instance(closer)
	defer func() {
		err := closer.Close()
		logger.Notify(err)
	}()

	upg, err := tableflip.New(tableflip.Options{})
	emperror.Panic(err)
	container.instance(upg)

	pipeline := new(run.Group)
	container.instance(pipeline)

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			logger.Info("graceful reloading")
			logger.Notify(upg.Upgrade())
		}
	}()
	ctx := context.Background()
	pipeline.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	pipeline.Add(appkitrun.GracefulRestart(ctx, upg))

	container.
		// base providers
		register(propagation.ProvideHTTP).
		register(metrics.Provide).
		register(httpserver.Provide).
		register(grpcProvider.Provide).
		// events
		register(eventrouter.Provide).
		register(events.ProvideKafkaPublisher).
		register(events.ProvideKafkaSubscriber).
		register(events.ProvideNatsSubscriber).
		register(events.ProvideNatsPublisher).
		register(events.ProvideChannel).
		// constructing application
		make(app).
		// boot components
		execute(jaegerexporter.Configure).
		execute(prometheus.Configure)

	// running application
	if err := pipeline.Run(); err != nil {
		logger.WithErrFilter(match.As(&run.SignalError{}).MatchError).Notify(err)
	}
}
