package prometheus

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/version"
)

type argsIn struct {
	dig.In
	Router  *http.ServeMux
	App     *config.Application
	Log     multilog.Logger
	Version version.Info
}

func Configure(in argsIn) error {
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
