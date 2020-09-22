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
	Log    *log.Config      `optional:"true"`
	Sentry *sentry.Config   `optional:"true"`
	Logger logur.Logger     `optional:"true"`
}

func Provide(in args) multilog.Logger {
	var opts []multilog.Option

	if in.Sentry != nil {
		opts = append(opts, multilog.WithHandler(sentry.Handler(in.Sentry.Addr, in.Sentry.Level, in.Sentry.Stop)))
	}

	if in.Log != nil {
		if in.Logger == nil {
			var writer io.Writer
			if in.Log.Color {
				zeroWriter := zerolog.NewConsoleWriter()
				zeroWriter.FormatMessage = func(i interface{}) string {
					return fmt.Sprintf("%-30s|", i)
				}
				writer = zeroWriter
			} else {
				writer = os.Stderr
			}

			logger := zerolog.New(writer).With().
				Timestamp().
				Logger()
			in.Logger = zerologadapter.New(logger)
		}

		opts = append(opts, multilog.WithHandler(log.Handler(in.Logger, in.Log.Stop)))
	}

	if in.Config == nil {
		in.Config = &multilog.Config{Level: logur.Error}
	}

	return multilog.New(in.Config.Level, opts...)
}
