package mongodb

import (
	"context"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/lebrains/gomongowrapper"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
)

type Args struct {
	dig.In
	Cfg    *Config
	Closer *shutdown.CloseGroup `optional:"true"`
	Health health.Health        `optional:"true"`
}

func Provide(in Args) (*gomongowrapper.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), in.Cfg.ConnectTimeout)
	defer cancel()

	client, err := gomongowrapper.NewConnector(in.Cfg.Config)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	if in.Closer != nil {
		in.Closer.Add(shutdown.SimpleCloseFn(func() {
			_ = client.Disconnect(ctx)
		}))
	}

	if in.Health != nil {
		err = in.Health.RegisterCheck(
			checks.Must(checks.NewPingCheck("mongodb.check", client)),
			health.ExecutionPeriod(3*time.Second),
		)
	}

	return client.Database(in.Cfg.Config.Name), err
}
