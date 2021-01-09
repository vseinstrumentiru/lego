package server

import (
	"github.com/vseinstrumentiru/lego/v2/metrics"
	"github.com/vseinstrumentiru/lego/v2/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/metrics/propagation"
	"github.com/vseinstrumentiru/lego/v2/transport/events"
	"github.com/vseinstrumentiru/lego/v2/transport/events/eventrouter"
	"github.com/vseinstrumentiru/lego/v2/transport/grpc"
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpclient"
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpserver"
	"github.com/vseinstrumentiru/lego/v2/transport/mysql"
	"github.com/vseinstrumentiru/lego/v2/transport/nats"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"
	"github.com/vseinstrumentiru/lego/v2/transport/stan"
)

func providers(runtime *Runtime) []interface{} {
	if runtime.Is(optWithoutProviders) {
		return []interface{}{}
	}

	return []interface{}{
		// metric providers
		propagation.ProvideHTTP,
		newrelicexporter.Provide,
		metrics.ProvideHealthChecker,
		metrics.ProvideMonitoringServer,
		// http / grpc
		httpserver.Provide,
		httpclient.Provide,
		httpclient.ConstructorProvider,
		grpc.Provide,
		// database
		mysql.Provide,
		sql.Provide,
		// events
		nats.Provide,
		stan.Provide,
		eventrouter.Provide,
		events.ProvideKafkaPublisher,
		events.ProvideKafkaSubscriber,
		events.ProvideNatsSubscriber,
		events.ProvideNatsPublisher,
		events.ProvideChannel,
	}
}
