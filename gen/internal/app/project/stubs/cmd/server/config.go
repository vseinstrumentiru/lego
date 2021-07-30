package main

import (
	cfg "github.com/vseinstrumentiru/lego/v2/config"
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/multilog/log"
	"github.com/vseinstrumentiru/lego/v2/transport/http"
)

type config struct {
	cfg.Application `mapstructure:",squash"`
	Modules         struct {
		Log     multilog.Config
		Console log.Config
		// Sentry sentry.Config
		HTTP http.Config
		// GRPC grpc.Config
		Metrics metrics.Config
		// Jaeger exporters.Jaeger
		// NewRelic exporters.NewRelic
		// Opencensus exporters.Opencensus
		// Tracing tracing.Config
		// Kafka kafka.Config
		// MongoDB mongodb.Config
		// MySQL mysql.Config
	}
}
