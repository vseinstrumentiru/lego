package opencensusexporter

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/internal/metrics/propagation"
	"github.com/vseinstrumentiru/lego/metrics/exporters"
	"github.com/vseinstrumentiru/lego/multilog"
)

type argsIn struct {
	dig.In
	App         *config.Application
	Config      *exporters.Opencensus `optional:"true"`
	Log         multilog.Logger
	Propagation *propagation.HTTPFormatCollection
}

func Configure(in argsIn) error {
	if in.Config == nil {
		return nil
	}

	log := in.Log.WithFields(map[string]interface{}{"component": "exporter.opencensus"})

	exp, err := ocagent.NewExporter(append(
		in.Config.Options(),
		ocagent.WithServiceName(in.App.Name),
	)...)

	if err != nil {
		return err
	}

	trace.RegisterExporter(exp)
	view.RegisterExporter(exp)
	in.Propagation.Add(&tracecontext.HTTPFormat{})

	log.Info("opencensus enabled")

	return nil
}
