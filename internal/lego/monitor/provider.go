package monitor

import (
	"net/http"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"emperror.dev/emperror"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"

	lego2 "github.com/vseinstrumentiru/lego/internal/lego"
	lepropagation "github.com/vseinstrumentiru/lego/internal/lego/monitor/propagation"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	// jaegerPropagation "contrib.go.opencensus.io/exporter/jaeger/propagation"
	jaegerPropagation "github.com/vseinstrumentiru/lego/internal/lego/monitor/propagation/jaegerwrap"
	"github.com/vseinstrumentiru/lego/internal/lego/monitor/telemetry"
)

func Provide(p lego2.Process, config Config) (*http.ServeMux, health.Health) {
	router, healthz := telemetry.Provide(p, p.Build())

	var formatter lepropagation.HTTPFormatCollection

	if config.Trace.Enabled {
		trace.ApplyConfig(config.Trace.Config())
	}

	if config.Exporter.Opencensus.Enabled {
		p.Info("opencensus exporter enabled")

		exporter, err := ocagent.NewExporter(append(
			config.Exporter.OpencensusOptions(),
			ocagent.WithServiceName(p.Name()),
		)...)
		emperror.Panic(err)

		formatter = append(formatter, &tracecontext.HTTPFormat{})
		trace.RegisterExporter(exporter)
		view.RegisterExporter(exporter)
	}

	if config.Exporter.NewRelic.Enabled {
		p.Info("newrelic exporter enabled")

		exporter, err := nrcensus.NewExporter(
			p.Name(),
			config.Exporter.NewRelic.Key,
		)
		emperror.Panic(err)

		trace.RegisterExporter(exporter)
		view.RegisterExporter(exporter)
	}

	if config.Exporter.Jaeger.Enabled {
		p.Info("jaeger exporter enabled")

		exporter, err := jaeger.NewExporter(jaeger.Options{
			CollectorEndpoint: config.Exporter.Jaeger.Addr,
			Process: jaeger.Process{
				ServiceName: p.Name(),
			},
			OnError: p.Handle,
		})

		emperror.Panic(err)

		formatter = append(formatter, &jaegerPropagation.HTTPFormat{})
		trace.RegisterExporter(exporter)
	}

	if config.Exporter.Prometheus.Enabled {
		p.Info("prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(prometheus.Options{
			OnError: emperror.WithDetails(
				p.Handler(),
				"component", "opencensus",
				"exporter", "prometheus",
			).Handle,
			ConstLabels: map[string]string{
				"app": p.Name(),
			},
		})
		emperror.Panic(err)

		view.RegisterExporter(exporter)
		router.Handle("/metrics", exporter)
	}

	lepropagation.DefaultHTTPFormat = formatter

	return router, healthz
}
