package events

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	watermilllog "logur.dev/integration/watermill"

	"github.com/vseinstrumentiru/lego/multilog"
)

type channelArgs struct {
	Config *gochannel.Config `optional:"true"`
	Logger multilog.Logger
}

func ProvideChannel(in channelArgs) (*gochannel.GoChannel, error) {
	cfg := gochannel.Config{}

	if in.Config != nil {
		cfg = *in.Config
	}

	ch := gochannel.NewGoChannel(
		cfg,
		watermilllog.New(in.Logger.WithFields(map[string]interface{}{"component": "events.pubsub", "provider": "gochannel"})),
	)

	return ch, nil
}
