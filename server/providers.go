package server

import (
	"github.com/vseinstrumentiru/lego/v2/internal/events"
	"github.com/vseinstrumentiru/lego/v2/internal/events/eventrouter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/exporters/newrelicexporter"
	"github.com/vseinstrumentiru/lego/v2/internal/metrics/propagation"
	grpcProvider "github.com/vseinstrumentiru/lego/v2/internal/transpoort/grpc"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/http/httpclient"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/http/httpserver"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/mysql"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/nats"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/sql"
	"github.com/vseinstrumentiru/lego/v2/internal/transpoort/stan"
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
		grpcProvider.Provide,
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
