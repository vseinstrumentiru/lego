package exporter

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/pkg/lego"
	"time"
)

type Config struct {
	Jaeger struct {
		lego.WithSwitch `mapstructure:",squash"`
		Addr            string
	}

	Opencensus struct {
		lego.WithSwitch `mapstructure:",squash"`

		Addr            string
		Insecure        bool
		ReconnectPeriod time.Duration
	}

	Prometheus struct {
		lego.WithSwitch `mapstructure:",squash"`
	}
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.monitor.exporter.prometheus.enabled", true)
}

func (c Config) Validate() (err error) {
	if c.Jaeger.Enabled && c.Jaeger.Addr == "" {
		err = errors.Append(err, errors.New("srv.monitor.exporter.jaeger.addr is required"))
	}

	if c.Opencensus.Enabled {
		if c.Opencensus.Addr == "" {
			err = errors.Append(err, errors.New("srv.monitor.exporter.opencensus.addr is required"))
		}
	}

	return
}

func (c Config) OpencensusOptions() []ocagent.ExporterOption {
	options := []ocagent.ExporterOption{
		ocagent.WithAddress(c.Opencensus.Addr),
		ocagent.WithReconnectionPeriod(c.Opencensus.ReconnectPeriod),
	}

	if c.Opencensus.Insecure {
		options = append(options, ocagent.WithInsecure())
	}

	return options
}
