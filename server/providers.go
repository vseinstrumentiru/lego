package server

import (
	"github.com/vseinstrumentiru/lego/v2/internal/metrics"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/propagation"
	events2 "github.com/vseinstrumentiru/lego/v2/transport/events"
	eventrouter2 "github.com/vseinstrumentiru/lego/v2/transport/events/eventrouter"
	"github.com/vseinstrumentiru/lego/v2/transport/grpc"
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpclient"
	"github.com/vseinstrumentiru/lego/v2/transport/http/httpserver"
	mysql2 "github.com/vseinstrumentiru/lego/v2/transport/mysql"
	nats2 "github.com/vseinstrumentiru/lego/v2/transport/nats"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"
	stan2 "github.com/vseinstrumentiru/lego/v2/transport/stan"
)

func providers() []interface{} {
	return []interface{}{
		// metric providers
		propagation.ProvideHTTP,
		newrelicexporter.Provide,
		metrics.Provide,
		// http / grpc
		httpserver.Provide,
		httpclient.Provide,
		httpclient.ConstructorProvider,
		grpc.Provide,
		// database
		mysql2.Provide,
		sql.Provide,
		// events
		nats2.Provide,
		stan2.Provide,
		eventrouter2.Provide,
		events2.ProvideKafkaPublisher,
		events2.ProvideKafkaSubscriber,
		events2.ProvideNatsSubscriber,
		events2.ProvideNatsPublisher,
		events2.ProvideChannel,
	}
}
