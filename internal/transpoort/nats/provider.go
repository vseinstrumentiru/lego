package nats

import (
	"emperror.dev/errors"
	"github.com/nats-io/nats.go"
	"go.uber.org/dig"

	lenats "github.com/vseinstrumentiru/lego/transport/nats"
	lestan "github.com/vseinstrumentiru/lego/transport/stan"
)

type args struct {
	dig.In
	Config     *lenats.Config
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

	if in.Config.AllowReconnect != nil && !*in.Config.AllowReconnect {
		options = append(options, nats.NoReconnect())
	}

	if in.Config.MaxReconnect != nil {
		options = append(options, nats.MaxReconnects(*in.Config.MaxReconnect))
	} else {
		options = append(options, nats.MaxReconnects(-1))
	}

	if in.Config.ReconnectWait != nil {
		options = append(options, nats.ReconnectWait(*in.Config.ReconnectWait))
	}

	options = append(options, nats.ReconnectBufSize(-1))

	conn, err := nats.Connect(in.Config.Addr, options...)

	return conn, err
}
