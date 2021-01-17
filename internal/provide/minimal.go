package provide

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
)

func Minimal() []interface{} {
	return []interface{}{
		propagation.ProvideHTTP,
		metrics.ProvideHealthChecker,
	}
}
