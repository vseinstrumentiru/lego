package nats

import (
	"emperror.dev/errors"
	"github.com/nats-io/nats.go"
	"go.uber.org/dig"

	lestan "github.com/vseinstrumentiru/lego/v2/transport/stan"
)

type args struct {
	dig.In
	Config     *Config
	StanConfig *lestan.Config `optional:"true"`
}

func Provide(in args) (*nats.Conn, error) {
	if in.Config == nil {
		return nil, errors.New("nats config not found")
	}

	var options []nats.Option

	if in.Config.Name != "" {
		options = append(options, nats.Name(in.Config.Name))
	} else if in.StanConfig != nil {
		options = append(options, nats.Name(in.StanConfig.GetClientID()))
	}

	if !in.Config.AllowReconnect {
		options = append(options, nats.NoReconnect())
	}

	options = append(options,
		nats.MaxReconnects(in.Config.MaxReconnect),
		nats.ReconnectWait(in.Config.ReconnectWait),
		nats.ReconnectBufSize(-1),
	)

	conn, err := nats.Connect(in.Config.Addr, options...)

	return conn, err
}
