package jaegerexporter

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
)

type Args struct {
	dig.In
	App         *config.Application
	Config      *exporters.Jaeger `optional:"true"`
	Log         log.Logger
	Propagation *propagation.HTTPFormatCollection `optional:"true"`
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
	if in.Propagation != nil {
		in.Propagation.Add(&HTTPFormat{})
	}

	log.Info("jaeger enabled")

	return nil
}
