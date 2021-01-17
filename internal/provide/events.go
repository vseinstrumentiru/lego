package provide

import (
	"github.com/vseinstrumentiru/lego/v2/transport/events"
	"github.com/vseinstrumentiru/lego/v2/transport/events/eventrouter"
	"github.com/vseinstrumentiru/lego/v2/transport/nats"
	"github.com/vseinstrumentiru/lego/v2/transport/stan"
)

func Events() []interface{} {
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
