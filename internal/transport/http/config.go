package http

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
	flag.String("http-addr", ":8000", "HTTP server address")
	_ = env.BindPFlag("srv.http.addr", flag.Lookup("http-addr"))
	env.SetDefault("srv.http.addr", ":8080")
}

func (c Config) Validate() (err error) {
	if c.Enabled {
		if c.Addr == "" {
			err = errors.Append(err, errors.New("srv.http.addr is required"))
		}
	}

	return
}
