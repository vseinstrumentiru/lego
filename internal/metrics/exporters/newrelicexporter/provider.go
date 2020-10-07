package newrelicexporter

import (
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/metrics/exporters"
	"github.com/vseinstrumentiru/lego/multilog"
)

type args struct {
	dig.In
	App    *config.Application
	Config *exporters.NewRelic `optional:"true"`
	Logger multilog.Logger
}

type loggerWrap struct {
	multilog.Logger
}

func (l loggerWrap) Error(msg string, context map[string]interface{}) {
	l.Logger.Error(msg, context)
}

func (l loggerWrap) Warn(msg string, context map[string]interface{}) {
	l.Logger.Warn(msg, context)
}

func (l loggerWrap) Info(msg string, context map[string]interface{}) {
	l.Logger.Info(msg, context)
}

func (l loggerWrap) Debug(msg string, context map[string]interface{}) {
	l.Logger.Debug(msg, context)
}

func (l loggerWrap) DebugEnabled() bool {
	return false
}

func Provide(in args) (app *newrelic.Application, err error) {
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
