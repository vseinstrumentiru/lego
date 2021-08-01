package prometheus

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type Args struct {
	dig.In
	Router  *http.ServeMux `optional:"true"`
	App     *config.Application
	Log     log.Logger
	Version *version.Info
}

func Configure(in Args) error {
	if in.Router == nil {
		return nil
	}

	log := in.Log.WithFields(map[string]interface{}{"component": "exporter.prometheus"})

	exp, err := prometheus.NewExporter(prometheus.Options{
		OnError: log.Handle,
		ConstLabels: map[string]string{
			"app":  in.App.Name,
			"host": in.Version.Host,
		},
	})
	if err != nil {
		return err
	}

	view.RegisterExporter(exp)
	in.Router.Handle("/metrics", exp)

	log.Info("prometheus enabled")

	return nil
}
