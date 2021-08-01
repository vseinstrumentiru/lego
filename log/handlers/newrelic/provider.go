package newrelic

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/log/handlers"
)

type Args struct {
	dig.In
	Newreilc *newrelic.Application `optional:"true"`
}

func Provide(in Args) handlers.Out {
	if in.Newreilc == nil {
		return handlers.Empty()
	}

	return handlers.Provider(NewHandler(in.Newreilc))
}
