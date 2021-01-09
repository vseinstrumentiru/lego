package server

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/opencensusexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/prometheus"
)

func executors(r *Runtime) (exec []interface{}) {
	if r.Is(optWithoutProviders) {
		return nil
	}

	return []interface{}{
		jaegerexporter.Configure,
		prometheus.Configure,
		opencensusexporter.Configure,
		newrelicexporter.Configure,
		metrics.ConfigureTrace,
		metrics.ConfigureStats,
	}
}
