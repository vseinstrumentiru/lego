package grpc

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/pkg/lego"
)

type Config struct {
	lego.WithSwitch `mapstructure:",squash"`

	Addr string
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	flag.String("grpc-addr", ":8001", "GRPC server address")
	_ = env.BindPFlag("srv.grpc.addr", flag.Lookup("grpc-addr"))
	env.SetDefault("srv.grpc.addr", ":8001")
}

func (c Config) Validate() (err error) {
	if c.Enabled {
		if c.Addr == "" {
			err = errors.Append(err, errors.New("srv.grpc.addr is required"))
		}
	}

	return
}
