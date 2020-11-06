package stan

import (
	"emperror.dev/errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"go.uber.org/dig"

	lestan "github.com/vseinstrumentiru/lego/v2/transport/stan"
)

type args struct {
	dig.In
	Stan *lestan.Config
	Nats *nats.Conn
}

func Provide(in args) (stan.Conn, error) {
	if in.Stan == nil {
		return nil, errors.New("stan config not found")
	}
	var options []stan.Option

	if in.Nats != nil {
		options = append(options, stan.NatsConn(in.Nats))
	} else {
		return nil, errors.New("nats connect not found")
	}

	if in.Stan.AckTimeout != nil {
		options = append(options, stan.PubAckWait(*in.Stan.AckTimeout))
	}

	options = append(options,
		stan.ConnectWait(in.Stan.ConnectTimeout),
		stan.MaxPubAcksInflight(in.Stan.MaxPubAcksInflight),
		stan.Pings(in.Stan.PingInterval, in.Stan.PingMaxOut),
	)

	conn, err := stan.Connect(in.Stan.ClusterID, in.Stan.GetClientID(), options...)

	return conn, err
}
