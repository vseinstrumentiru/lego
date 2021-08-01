package newrelicexporter

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/log"
	lenewrelic "github.com/vseinstrumentiru/lego/v2/log/handlers/newrelic"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters"
)

type ConfigArgs struct {
	dig.In
	App      *config.Application   `optional:"true"`
	Config   *exporters.NewRelic   `optional:"true"`
	NewRelic *newrelic.Application `optional:"true"`
	Log      log.Logger
}

func Configure(in ConfigArgs) error {
	if in.Config == nil || !in.Config.Enabled {
		return nil
	}

	in.Log.WithHandler(lenewrelic.NewHandler(in.NewRelic))

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
