package module

import "github.com/vseinstrumentiru/lego/v2/di"

func All() []di.Module {
	return []di.Module{
		Pipeline,
		Upgrader,
		SQLDB,
		HealthChecker,
		NATS,
		STAN,
		NATSPublisher,
		NATSSubscriber,
		KafkaPublisher,
		KafkaSubscriber,
		EventRouter,
		HTTPServer,
		HTTPClient,
		HTTPClientConstructor,
		MuxRouter,
		GRPCServer,
		MonitoringServer,
		HTTPPropagation,
		NewRelicExporter,
	}
}
