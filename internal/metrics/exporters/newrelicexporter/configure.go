package newrelicexporter

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/metrics/exporters"
	"github.com/vseinstrumentiru/lego/multilog"
	lenewrelic "github.com/vseinstrumentiru/lego/multilog/newrelic"
)

type argsIn struct {
	dig.In
	App      *config.Application   `optional:"true"`
	Config   *exporters.NewRelic   `optional:"true"`
	NewRelic *newrelic.Application `optional:"true"`
	Log      multilog.Logger
}

func Configure(in argsIn) error {
	if in.Config == nil || !in.Config.Enabled {
		return nil
	}

	in.Log.WithHandler(lenewrelic.Handler(in.NewRelic))

	if !in.Config.TelemetryEnabled {
		return nil
	}

	log := in.Log.WithFields(map[string]interface{}{"component": "exporter.newrelic"})

	exp, err := nrcensus.NewExporter(in.App.FullName(), in.Config.Key)

	if err != nil {
		return err
	}

	trace.RegisterExporter(exp)
	view.RegisterExporter(exp)

	log.Info("newrelic enabled")

	return nil
}
