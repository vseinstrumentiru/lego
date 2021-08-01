package file

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/vseinstrumentiru/lego/v2/log/handlers"
	"github.com/vseinstrumentiru/lego/v2/log/handlers/console"
)

type Args struct {
	Config *Config `optional:"true"`
}

func Provide(in Args) handlers.Out {
	if in.Config == nil {
		return handlers.Empty()
	}

	zeroWriter := zerolog.NewConsoleWriter()

	zeroWriter.TimeFormat = console.DefaultTimeFormat
	zeroWriter.NoColor = true

	zeroWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%-30s|", i)
	}

	zeroWriter.Out = &in.Config.Logger

	ctx := zerolog.New(zeroWriter).With().
		Timestamp()

	logger := console.NewAdapter(ctx.Logger().Level(console.ToZerologLevel(in.Config.Level)))

	return handlers.Provider(console.NewHandler(logger, in.Config.Stop))
}
