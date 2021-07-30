package multilogprovider

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/multilog/log"
	"github.com/vseinstrumentiru/lego/v2/multilog/sentry"
)

type args struct {
	dig.In
	App    *config.Application
	Config *multilog.Config `optional:"true"`

	Sentry *sentry.Config `optional:"true"`
	Log    *log.Config    `optional:"true"`
	Logger logur.Logger   `optional:"true"`
}

type ctxOpt func(ctx zerolog.Context) zerolog.Context

func withCaller(depth int) ctxOpt {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.CallerWithSkipFrameCount(depth)
	}
}

func Provide(in args) multilog.Logger {
	var opts []multilog.Option

	if in.Config == nil {
		level := logur.Error
		if in.App.DebugMode {
			level = logur.Trace
		}

		in.Config = &multilog.Config{Level: level}
	}

	if in.Sentry != nil {
		opts = append(opts, multilog.WithHandler(sentry.Handler(in.Sentry.Addr, in.Sentry.Level, in.Sentry.Stop)))
	}

	var contextOptions []ctxOpt

	if in.Log == nil && !in.Config.SilentMode {
		in.Log = log.DefaultConfig()

		if in.App.LocalMode {
			in.Log.Color = true
		}
	}

	//nolint:nestif
	if in.Log != nil {
		if in.Logger == nil {
			var writer io.Writer
			if in.Log.TimeFormat != "" {
				zerolog.TimeFieldFormat = in.Log.TimeFormat
			}

			if in.Log.Color {
				zeroWriter := zerolog.NewConsoleWriter()
				if in.Log.TimeFormat != "" {
					zeroWriter.TimeFormat = in.Log.TimeFormat
				}
				zeroWriter.FormatMessage = func(i interface{}) string {
					return fmt.Sprintf("%-30s|", i)
				}
				writer = zeroWriter
				if in.Log.Depth > 0 {
					contextOptions = append(contextOptions, withCaller(in.Log.Depth))
				}
			} else {
				writer = os.Stderr
			}

			ctx := zerolog.New(writer).With().
				Timestamp()

			for _, o := range contextOptions {
				ctx = o(ctx)
			}

			logger := ctx.Logger().Level(zerolog.Level(in.Log.Level - 1))

			in.Logger = NewAdapter(logger)
		}

		opts = append(opts, multilog.WithHandler(log.Handler(in.Logger, in.Log.Stop)))
	}

	logger := multilog.New(in.Config.Level, opts...)

	if in.App.Name != config.Undefined && in.App.Name != "" {
		fields := make(map[string]interface{})
		fields["app"] = in.App.Name

		logger = logger.WithFields(fields)
	}

	return logger
}
