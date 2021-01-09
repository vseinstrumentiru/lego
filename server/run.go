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

	baseCfg "github.com/vseinstrumentiru/lego/v2/config"
	di "github.com/vseinstrumentiru/lego/v2/internal/container"
	"github.com/vseinstrumentiru/lego/v2/internal/deprecated"
	environment "github.com/vseinstrumentiru/lego/v2/internal/env"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/multilog/multilogprovider"
	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type Option func(r *Runtime)

func Run(app interface{}, opts ...Option) {
	runtime := newRuntime(opts)
	// core container instance
	container := newContainer()

	env := runtime.newEnv()
	env.SetFlag("version", false, "show version")
	var ver version.Info

	container.
		register(func() di.Container { return container.i }).
		register(func() environment.Env { return env })
	// configuration
	runtime.onConfig(func(cfg interface{}) {
		container.register(func() environment.Config { return cfg })
	})

	container.make(env)

	if runtime.Is(optLocalDebug) {
		container.execute(func(cfg *baseCfg.Application) {
			cfg.LocalMode = true
			cfg.DebugMode = true
		})
	}

	container.
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
	container.register(multilogprovider.Provide).
		execute(func(l multilog.Logger) {
			logger = l
		})

	logger = logger.WithFields(map[string]interface{}{"component": "runtime"})
	logger.Trace("starting application")

	defer emperror.HandleRecover(logger.WithFields(map[string]interface{}{"type": "panic"}))
	multilog.SetStandardLogger(logger.WithFields(map[string]interface{}{"type": "standard"}))

	// pipeline configuration
	logger.Trace("pipeline configuration")
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

	logger.Trace("pipeline configuration", map[string]interface{}{"status": "completed"})

	// constructing application
	logger.Trace("constructing application")
	deprecated.Container = container.i
	for _, provider := range providers(runtime) {
		container.register(provider)
	}

	container.make(app)

	for _, exec := range executors(runtime) {
		container.execute(exec)
	}
	logger.Trace("constructing application", map[string]interface{}{"status": "completed"})

	if runtime.Not(optNoWait) {
		logger.Trace("starting pipeline")
		// running application
		if err := pipeline.Run(); err != nil {
			logger.WithErrFilter(match.As(&run.SignalError{}).MatchError).Notify(err)
		}
	}
	logger.Trace("application stopped")
}
