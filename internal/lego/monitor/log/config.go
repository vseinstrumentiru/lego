package log

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Format    string
	Level     string
	NoColor   bool
	UseStack  bool
	SkipStack int
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.monitor.log.format", "json")
	env.SetDefault("srv.monitor.log.level", "debug")
	env.SetDefault("srv.monitor.log.skipStack", 4)

	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		env.SetDefault("no_color", true)
	}

	env.RegisterAlias("srv.log.noColor", "no_color")
}

func (c *Config) Validate() (err error) {
	return
}
