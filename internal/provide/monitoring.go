package provide

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
)

func Monitoring() []interface{} {
	return []interface{}{
		newrelicexporter.Provide,
		metrics.ProvideMonitoringServer,
	}
}
