package server

import (
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/opencensusexporter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/prometheus"
)

func executors() []interface{} {
	return []interface{}{
		jaegerexporter.Configure,
		prometheus.Configure,
		opencensusexporter.Configure,
		newrelicexporter.Configure,
	}
}
