package grpc

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/pkg/lego"
)

type Config struct {
	lego.WithSwitch `mapstructure:",squash"`

	Port int
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	flag.String("grpc-port", "8001", "GRPC server port")
	_ = env.BindPFlag("srv.grpc.port", flag.Lookup("grpc-addr"))
	env.SetDefault("srv.grpc.port", 8001)
}

func (c Config) Validate() (err error) {
	if c.Enabled {
		if c.Port == 0 {
			err = errors.Append(err, errors.New("srv.grpc.port is required"))
		}
	}

	return
}
