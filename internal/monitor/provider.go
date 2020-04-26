package monitor

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	jaegerPropagation "contrib.go.opencensus.io/exporter/jaeger/propagation"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"emperror.dev/emperror"
	health "github.com/AppsFlyer/go-sundheit"
	lepropagation "github.com/vseinstrumentiru/lego/internal/monitor/propagation"
	"github.com/vseinstrumentiru/lego/internal/monitor/telemetry"
	"github.com/vseinstrumentiru/lego/pkg/build"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"net/http"
)

func ProvideServer(p lego.Process, buildinfo build.Info, config Config) (*http.ServeMux, health.Health) {
	router, healthz := telemetry.Provide(p, buildinfo)

	var formatter lepropagation.HTTPFormatCollection

	if config.Trace.Enabled {
		trace.ApplyConfig(config.Trace.Config())
	}

	if config.Exporter.Opencensus.Enabled {
		exporter, err := ocagent.NewExporter(append(
			config.Exporter.OpencensusOptions(),
			ocagent.WithServiceName(p.Name()),
		)...)
		emperror.Panic(err)

		formatter = append(formatter, &tracecontext.HTTPFormat{})
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
			Namespace: p.Name(),
			OnError: emperror.WithDetails(
				p.Handler(),
				"component", "opencensus",
				"exporter", "prometheus",
			).Handle,
		})
		emperror.Panic(err)

		lepropagation.DefaultHTTPFormat = formatter
		view.RegisterExporter(exporter)
		router.Handle("/metrics", exporter)
	}

	return router, healthz
}

func Provide(p lego.Process, buildinfo build.Info, config Config) (*http.ServeMux, health.Health) {
	router, healthz := telemetry.Provide(p, buildinfo)

	if config.Trace.Enabled {
		trace.ApplyConfig(config.Trace.Config())
	}

	if config.Exporter.Opencensus.Enabled {
		exporter, err := ocagent.NewExporter(append(
			config.Exporter.OpencensusOptions(),
			ocagent.WithServiceName(p.Name()),
		)...)
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

		trace.RegisterExporter(exporter)
	}

	if config.Exporter.Prometheus.Enabled {
		p.Info("prometheus exporter enabled")

		exporter, err := prometheus.NewExporter(prometheus.Options{
			Namespace: p.Name(),
			OnError: emperror.WithDetails(
				p.Handler(),
				"component", "opencensus",
				"exporter", "prometheus",
			).Handle,
		})
		emperror.Panic(err)

		view.RegisterExporter(exporter)
		router.Handle("/metrics", exporter)
	}

	return router, healthz
}
