package multilog

import (
	"fmt"
	"io"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
	zerologadapter "logur.dev/adapter/zerolog"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/multilog"
	"github.com/vseinstrumentiru/lego/multilog/log"
	lenewrelic "github.com/vseinstrumentiru/lego/multilog/newrelic"
	"github.com/vseinstrumentiru/lego/multilog/sentry"
)

type args struct {
	dig.In
	Config *multilog.Config `optional:"true"`

	Sentry   *sentry.Config        `optional:"true"`
	NewRelic *newrelic.Application `optional:"true"`
	Log      *log.Config           `optional:"true"`
	Logger   logur.Logger          `optional:"true"`
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

	if in.NewRelic != nil {
		opts = append(opts, multilog.WithHandler(lenewrelic.Handler(in.NewRelic)))
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

			in.Logger = zerologadapter.New(logger)
		}

		opts = append(opts, multilog.WithHandler(log.Handler(in.Logger, in.Log.Stop)))
	}

	if in.Config == nil {
		in.Config = &multilog.Config{Level: logur.Error}
	}

	return multilog.New(in.Config.Level, opts...)
}
