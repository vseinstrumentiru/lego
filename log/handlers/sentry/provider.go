package sentry

import (
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/log/handlers"
)

type Args struct {
	dig.In
	Config *Config `optional:"true"`
}

func Provide(in Args) handlers.Out {
	if in.Config == nil {
		return handlers.Empty()
	}

	return handlers.Provider(NewHandler(
		in.Config.Addr,
		in.Config.Level,
		in.Config.Stop,
	))
}
