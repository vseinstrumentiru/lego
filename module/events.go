package module

import (
	"github.com/vseinstrumentiru/lego/v2/di"
	"github.com/vseinstrumentiru/lego/v2/transport/events"
	"github.com/vseinstrumentiru/lego/v2/transport/events/eventrouter"
	"github.com/vseinstrumentiru/lego/v2/transport/nats"
	"github.com/vseinstrumentiru/lego/v2/transport/stan"
)

func NATSPack() []di.Module {
	return []di.Module{
		NATS,
		STAN,
		NATSPublisher,
		NATSSubscriber,
	}
}

func NATSPublisherPack() []di.Module {
	return []di.Module{
		NATS,
		STAN,
		NATSPublisher,
	}
}

func NATSSubscriberPack() []di.Module {
	return []di.Module{
		NATS,
		STAN,
		NATSSubscriber,
	}
}

func KafkaPack() []di.Module {
	return []di.Module{
		KafkaPublisher,
		KafkaSubscriber,
	}
}

func EventRouter() (interface{}, []interface{}) {
	return eventrouter.Provide, nil
}

func NATS() (interface{}, []interface{}) {
	return nats.Provide, nil
}

func STAN() (interface{}, []interface{}) {
	return stan.Provide, nil
}

func NATSPublisher() (interface{}, []interface{}) {
	return events.ProvideNatsPublisher, nil
}

func NATSSubscriber() (interface{}, []interface{}) {
	return events.ProvideNatsSubscriber, nil
}

func ChannelPubSub() (interface{}, []interface{}) {
	return events.ProvideChannel, nil
}

func KafkaPublisher() (interface{}, []interface{}) {
	return events.ProvideKafkaPublisher, nil
}

func KafkaSubscriber() (interface{}, []interface{}) {
	return events.ProvideKafkaSubscriber, nil
}
