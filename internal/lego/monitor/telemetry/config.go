package telemetry

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"go.opencensus.io/stats/view"
)

type Config struct {
	// Telemetry HTTP server address
	Addr string

	Stats []*view.View
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	flag.String("telemetry-addr", ":10000", "Telemetry HTTP server address")
	_ = env.BindPFlag("srv.monitor.telemetry.addr", flag.Lookup("telemetry-addr"))
	env.SetDefault("srv.monitor.telemetry.addr", ":10000")
}

func (c Config) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("srv.monitor.telemetry.addr is required"))
	}

	return err
}
