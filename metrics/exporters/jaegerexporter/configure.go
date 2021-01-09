package jaegerexporter

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type Args struct {
	dig.In
	App         *config.Application
	Config      *exporters.Jaeger `optional:"true"`
	Log         multilog.Logger
	Propagation *propagation.HTTPFormatCollection
}

func Configure(in Args) error {
	if in.Config == nil {
		return nil
	}

	log := in.Log.WithFields(map[string]interface{}{"component": "exporter.jaeger"})

	exp, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: in.Config.Addr,
		Process: jaeger.Process{
			ServiceName: in.App.Name,
		},
		OnError: log.Handle,
	})

	if err != nil {
		return err
	}

	trace.RegisterExporter(exp)
	in.Propagation.Add(&HTTPFormat{})

	log.Info("jaeger enabled")

	return nil
}
