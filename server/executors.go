package server

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/jaegerexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/opencensusexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/prometheus"
)

func executors(r *Runtime) (exec []interface{}) {
	execList := map[string][]interface{}{
		optUseJaeger:     {jaegerexporter.Configure},
		optUsePrometheus: {prometheus.Configure},
		optUseOpencensus: {opencensusexporter.Configure},
		optUseNewRelic:   {newrelicexporter.Configure},
		optUseTrace:      {metrics.ConfigureTrace},
		optUseStats:      {metrics.ConfigureStats},
		optUseMonitoring: {metrics.ConfigureStats},
	}

	ignoreProviders := r.Is(optWithoutProviders)

	for key, cfg := range execList {
		if !ignoreProviders || r.Is(key) {
			exec = append(exec, cfg...)
		}
	}

	return exec
}
