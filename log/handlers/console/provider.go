package console

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/log/handlers"
)

type Args struct {
	dig.In
	Config *Config             `optional:"true"`
	Logger logur.Logger        `optional:"true"`
	App    *config.Application `optional:"true"`
}

type ctxOpt func(ctx zerolog.Context) zerolog.Context

func withCaller(depth int) ctxOpt {
	return func(ctx zerolog.Context) zerolog.Context {
		return ctx.CallerWithSkipFrameCount(depth)
	}
}

func Provide(in Args) handlers.Out {
	if in.Config == nil {
		in.Config = DefaultConfig()
	}

	if in.Config.Color {
		in.Config.Format = HumanFormat
	}

	if in.Config.out == nil {
		in.Config.out = os.Stderr
	}

	if in.App != nil {
		if in.App.LocalMode {
			in.Config.Format = HumanFormat
		}
	}

	var contextOptions []ctxOpt

	// nolint: nestif
	if in.Logger == nil {
		var writer io.Writer

		if in.Config.TimeFormat != "" {
			zerolog.TimeFieldFormat = in.Config.TimeFormat
		}

		if in.Config.Format == HumanFormat {
			zeroWriter := zerolog.NewConsoleWriter()
			if in.Config.TimeFormat != "" {
				zeroWriter.TimeFormat = in.Config.TimeFormat
			}

			zeroWriter.FormatMessage = func(i interface{}) string {
				return fmt.Sprintf("%-30s|", i)
			}

			zeroWriter.Out = in.Config.out

			writer = zeroWriter

			if in.Config.Depth > 0 {
				contextOptions = append(contextOptions, withCaller(in.Config.Depth))
			}
		} else {
			writer = in.Config.out
		}

		ctx := zerolog.New(writer).With().
			Timestamp()

		for _, o := range contextOptions {
			ctx = o(ctx)
		}

		logger := ctx.Logger().Level(ToZerologLevel(in.Config.Level))

		in.Logger = NewAdapter(logger)
	}

	return handlers.Provider(NewHandler(in.Logger, in.Config.Stop))
}
