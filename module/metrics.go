package module

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/opencensusexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/prometheus"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
)

func MonitoringServer() (interface{}, []interface{}) {
	return metrics.ProvideMonitoringServer, []interface{}{
		prometheus.Configure,
		jaegerexporter.Configure,
		opencensusexporter.Configure,
		newrelicexporter.Configure,
		metrics.ConfigureTrace,
		metrics.ConfigureStats,
	}
}

func HealthChecker() (interface{}, []interface{}) {
	return metrics.ProvideHealthChecker, nil
}

func HTTPPropagation() (interface{}, []interface{}) {
	return propagation.ProvideHTTP, nil
}

func NewRelicExporter() (interface{}, []interface{}) {
	return newrelicexporter.Provide, nil
}
