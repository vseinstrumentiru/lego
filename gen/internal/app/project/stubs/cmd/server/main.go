package main

import (
	lego "github.com/vseinstrumentiru/lego/v2/app"
	"github.com/vseinstrumentiru/lego/v2/module"
)

func main() {
	lego.NewRuntime(
		lego.ServerMode(),
		lego.WithConfig(&config{}),
		lego.Provide(
			module.Pipeline,
			module.Upgrader,
			module.MonitoringServer,
			// module.MySQLPack,
			// module.PostgresPack,
			// module.MongoDbPack,
			// module.KafkaPublisher,
			// module.KafkaSubscriber,
			module.HTTPServer,
			module.MuxRouter,
			// module.HTTPClient,
			// module.HTTPClientConstructor,
		),
	).Run(app{})
}
