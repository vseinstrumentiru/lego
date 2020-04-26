package monitor

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vseinstrumentiru/lego/internal/monitor/errorhandler"
	"github.com/vseinstrumentiru/lego/internal/monitor/exporter"
	"github.com/vseinstrumentiru/lego/internal/monitor/log"
	"github.com/vseinstrumentiru/lego/internal/monitor/telemetry"
	"github.com/vseinstrumentiru/lego/internal/monitor/tracer"
)

type Config struct {
	Log          log.Config
	Errorhandler errorhandler.Config
	Exporter     exporter.Config
	Trace        tracer.Config
	Telemetry    telemetry.Config
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	c.Log.SetDefaults(env, flag)
	c.Errorhandler.SetDefaults(env, flag)
	c.Exporter.SetDefaults(env, flag)
	c.Trace.SetDefaults(env, flag)
	c.Telemetry.SetDefaults(env, flag)
}

func (c Config) Validate() (err error) {
	err = errors.Append(err, c.Log.Validate())
	err = errors.Append(err, c.Errorhandler.Validate())
	err = errors.Append(err, c.Exporter.Validate())
	err = errors.Append(err, c.Trace.Validate())
	err = errors.Append(err, c.Telemetry.Validate())

	return
}
