package config

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/internal/monitor"
	"github.com/vseinstrumentiru/lego/internal/transport/event"
	"github.com/vseinstrumentiru/lego/internal/transport/grpc"
	"github.com/vseinstrumentiru/lego/internal/transport/http"
	"github.com/vseinstrumentiru/lego/pkg/build"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"time"
)

type WithSwitch = lego.WithSwitch

type Config struct {
	Name            string
	Env             string
	Debug           bool
	Host            string
	Http            http.Config
	Grpc            grpc.Config
	Events          event.Config
	Monitor         monitor.Config
	ShutdownTimeout time.Duration

	Build  build.Info
	Custom lego.Config `mapstructure:"-"`
}

func (c Config) Validate() (err error) {
	if c.Name == "" {
		err = errors.Append(err, errors.New("srv.name required"))
	}

	err = errors.Append(err, c.Http.Validate())
	err = errors.Append(err, c.Grpc.Validate())
	err = errors.Append(err, c.Events.Validate())
	err = errors.Append(err, c.Monitor.Validate())

	if c.Custom != nil {
		err = errors.Append(err, c.Custom.Validate())
	}

	return
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.env", "dev")
	env.SetDefault("srv.debug", false)
	env.SetDefault("srv.host", "localhost")
	env.SetDefault("srv.shutdownTimeout", 15*time.Second)

	c.Http.SetDefaults(env, flag)
	c.Grpc.SetDefaults(env, flag)
	c.Events.SetDefaults(env, flag)
	c.Monitor.SetDefaults(env, flag)
}
