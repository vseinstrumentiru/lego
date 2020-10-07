package multilog

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/dig"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/multilog"
	"github.com/vseinstrumentiru/lego/multilog/log"
	"github.com/vseinstrumentiru/lego/multilog/sentry"
)

type args struct {
	dig.In
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

	if in.Sentry != nil {
		opts = append(opts, multilog.WithHandler(sentry.Handler(in.Sentry.Addr, in.Sentry.Level, in.Sentry.Stop)))
	}

	var contextOptions []ctxOpt

	if in.Log != nil {
		if in.Logger == nil {
			var writer io.Writer
			if in.Log.Color {
				zeroWriter := zerolog.NewConsoleWriter()
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

			logger := ctx.Logger()

			if in.Log.Level >= 1 {
				logger = logger.Level(zerolog.Level(in.Log.Level - 1))
			} else {
				logger = logger.Level(zerolog.Level(in.Log.Level))
			}

			in.Logger = zerologadapter.New(logger)
		}

		opts = append(opts, multilog.WithHandler(log.Handler(in.Logger, in.Log.Stop)))
	}

	if in.Config == nil {
		in.Config = &multilog.Config{Level: logur.Error}
	}

	return multilog.New(in.Config.Level, opts...)
}
