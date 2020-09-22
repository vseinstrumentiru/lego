package events

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/events"
	"github.com/vseinstrumentiru/lego/multilog"
)

type channelArgs struct {
	Config *gochannel.Config `optional:"true"`
	Logger multilog.Logger
}

func ProvideChannel(in channelArgs) (events.ChannelProvider, error) {
	cfg := gochannel.Config{}

	if in.Config != nil {
		cfg = *in.Config
	}

	ch := gochannel.NewGoChannel(
		cfg,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"component": "events.pubsub", "provider": "gochannel"})),
	)

	return channelProvider{ch}, nil
}

type channelProvider struct {
	*gochannel.GoChannel
}

func (c channelProvider) NewPublisher() (message.Publisher, error) {
	return c, nil
}

func (c channelProvider) NewSubscriber() (message.Subscriber, error) {
	return c, nil
}
