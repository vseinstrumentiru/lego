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

func provideMinimal() []interface{} {
	return []interface{}{
		propagation.ProvideHTTP,
		metrics.ProvideHealthChecker,
	}
}

func provideMonitoring() []interface{} {
	return []interface{}{
		newrelicexporter.Provide,
		metrics.ProvideMonitoringServer,
	}
}

func provideHttp() []interface{} {
	return []interface{}{
		httpserver.Provide,
		httpclient.Provide,
		httpclient.ConstructorProvider,
	}
}

func provideGrpc() []interface{} {
	return []interface{}{
		grpc.Provide,
	}
}

func provideSql() []interface{} {
	return []interface{}{
		mysql.ProvideConnector,
		sql.Provide,
	}
}

func provideEvents() []interface{} {
	return []interface{}{
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

func providers(runtime *Runtime) []interface{} {
	if runtime.Is(optWithoutProviders) {
		return []interface{}{}
	}

	var res []interface{}

	res = append(res, provideMinimal()...)

	if runtime.Not(optWithoutMonitoring) {
		res = append(res, provideMonitoring()...)
	}

	res = append(res, provideHttp()...)
	res = append(res, provideGrpc()...)
	res = append(res, provideSql()...)
	res = append(res, provideEvents()...)

	return res
}
