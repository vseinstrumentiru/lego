package newrelicexporter

import (
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters"
	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type ProvideArgs struct {
	dig.In
	App    *config.Application
	Config *exporters.NewRelic `optional:"true"`
	Logger multilog.Logger
}

func Provide(in ProvideArgs) (app *newrelic.Application, err error) {
	if in.Config == nil {
		return nil, nil
	}

	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName(in.App.FullName()),
		newrelic.ConfigLicense(in.Config.Key),
		newrelic.ConfigLogger(loggerWrap{Logger: in.Logger}),
		newrelic.ConfigEnabled(in.Config.Enabled),
	)

	return
}
