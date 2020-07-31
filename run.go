package LeGo

import (
	"context"
	"emperror.dev/emperror"
	"emperror.dev/errors/match"
	"errors"
	"github.com/oklog/run"
	appkiterrors "github.com/sagikazarmark/appkit/errors"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/telemetry"
	"github.com/vseinstrumentiru/lego/internal/lego/transport/event"
	"github.com/vseinstrumentiru/lego/internal/lego/transport/grpc"
	"github.com/vseinstrumentiru/lego/internal/lego/transport/http"
	"github.com/vseinstrumentiru/lego/tools/contexttool"
	"logur.dev/logur"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context, app lego.App) {
	s := newServer(app.GetName(), extractConfig(app))
	provideConfig(app, s.Config.Custom)

	defer emperror.HandleRecover(s.Handler())

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			s.Info("graceful reloading")

			_ = s.Upgrader.Upgrade()
		}
	}()

	{
		config := s.Config.Monitor.Telemetry

		if cApp, ok := app.(lego.AppWithStats); ok {
			config.Stats = cApp.GetStats()
		}

		closer := telemetry.Run(s, s.Telemetry, config)
		defer closer.Close()
	}

	{
		const name = "app"
		app.SetLogErr(
			lego.NewContextAwareLogErr(
				logur.WithField(s.Log(), "server", name),
				s.Handler(),
				lecontext.Extractor,
			).
				WithFilter(appkiterrors.IsServiceError),
		)
	}

	pubApp, pubOk := app.(lego.AppWithPublishers)
	subApp, subOk := app.(lego.AppWithEventHandlers)
	pubOk, subOk = pubOk && s.Config.Events.Enabled, subOk && s.Config.Events.Enabled

	if pubOk || subOk {
		em, exec, interrupt := event.Run(s, s.Config.Events)
		defer em.Close()

		if pubOk {
			err := pubApp.RegisterEventDispatcher(em.Publisher())
			emperror.Panic(err)
		}

		if cApp, ok := app.(lego.AppWithRegistration); ok {
			closer, err := cApp.Register(s)
			emperror.Panic(err)

			defer lego.Close(closer)
		}

		if subOk {
			err := subApp.RegisterEventHandlers(em)
			emperror.Panic(err)
			s.Background(exec, interrupt)
		}
	} else if cApp, ok := app.(lego.AppWithRegistration); ok {
		closer, err := cApp.Register(s)
		emperror.Panic(err)

		defer lego.Close(closer)
	}

	if httpApp, ok := app.(lego.AppWithHttp); ok {
		if !s.Config.Http.Enabled {
			emperror.Panic(errors.New("http config not defined"))
		}

		httpRouter, closer := http.Run(s, s.Config.Http)
		defer closer.Close()

		err := httpApp.RegisterHTTP(httpRouter)
		emperror.Panic(err)
	}

	if grpcApp, ok := app.(lego.AppWithGrpc); ok {
		if !s.Config.Grpc.Enabled {
			emperror.Panic(errors.New("grpc config not defined"))
		}

		server, closer := grpc.Run(s, s.Config.Grpc)
		defer closer.Close()

		err := grpcApp.RegisterGRPC(server)
		emperror.Panic(err)
	}

	if runApp, ok := app.(lego.AppWithRunner); ok {
		terminate := make(chan bool, 1)
		s.Runner.Add(func() error {
			return runApp.Run(terminate)
		}, func(err error) {
			terminate <- true
		})
	}

	// Setup signal handler
	s.Runner.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	// Setup graceful restart
	s.Runner.Add(appkitrun.GracefulRestart(ctx, s.Upgrader))

	err := s.Runner.Run()
	if err != nil {
		s.WithFilter(match.As(&run.SignalError{}).MatchError).Handle(err)
	}
}
