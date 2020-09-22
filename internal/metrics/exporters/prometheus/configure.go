package prometheus

import (
	"net/http"
	"strings"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/config"
	"github.com/vseinstrumentiru/lego/multilog"
)

type argsIn struct {
	dig.In
	Router *http.ServeMux
	App    *config.Application
	Log    multilog.Logger
}

func Configure(in argsIn) error {
	log := in.Log.WithFields(map[string]interface{}{"component": "exporter.prometheus"})

	exp, err := prometheus.NewExporter(prometheus.Options{
		Namespace: strings.ReplaceAll(in.App.Name, "-", "_"),
		OnError:   log.Handle,
	})

	if err != nil {
		return err
	}

	view.RegisterExporter(exp)
	in.Router.Handle("/metrics", exp)

	log.Info("prometheus enabled")

	return nil
}
