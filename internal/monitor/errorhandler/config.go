package errorhandler

import (
	"emperror.dev/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Providers []string
	Sentry    struct {
		DSN string
	}
}

const (
	LogProvider    = "log"
	SentryProvider = "sentry"
)

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	env.SetDefault("srv.monitor.errorhandler.providers", []string{LogProvider})
}

func (c Config) Validate() (err error) {
	providers := map[string]bool{
		LogProvider:    false,
		SentryProvider: false,
	}

	for _, p := range c.Providers {
		if used, ok := providers[p]; !ok {
			err = errors.Append(err, errors.Errorf("errorHandler provider %s not defined", p))
		} else if used {
			err = errors.Append(err, errors.Errorf("errorHandler provider %s must be defined only once", p))
		}

		providers[p] = true
	}

	if providers[SentryProvider] && c.Sentry.DSN == "" {
		err = errors.Append(err, errors.New("sentry errorHandler provider requires uri"))
	}

	return err
}
