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

type Server struct {
	Name            string
	Env             string
	Debug           bool
	Http            http.Config
	Grpc            grpc.Config
	Events          event.Config
	Monitor         monitor.Config
	ShutdownTimeout time.Duration

	Build build.Info
	App   lego.Config `mapstructure:"-"`
}

func (c Server) Validate() (err error) {
	if c.Name == "" {
		err = errors.Append(err, errors.New("srv.name required"))
	}

	err = errors.Append(err, c.Http.Validate())
	err = errors.Append(err, c.Grpc.Validate())
	err = errors.Append(err, c.Events.Validate())
	err = errors.Append(err, c.Monitor.Validate())

	return
}

func (c Server) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.env", "dev")
	env.SetDefault("srv.debug", false)
	env.SetDefault("srv.shutdownTimeout", 15*time.Second)

	c.Http.SetDefaults(env, flag)
	c.Grpc.SetDefaults(env, flag)
	c.Events.SetDefaults(env, flag)
	c.Monitor.SetDefaults(env, flag)
}
