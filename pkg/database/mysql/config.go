package mysql

import (
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	mysql.Config `mapstructure:",squash"`
}

func (c Config) SetDefaults(key string, env *viper.Viper, flag *pflag.FlagSet) {
	key = strings.TrimSuffix(key, ".")

	if key == "" {
		emperror.Panic(errors.New("config key is empty"))
	}

	env.SetDefault(key+".parseTime", true)
	env.SetDefault(key+".rejectReadOnly", true)
	env.RegisterAlias(key+".passwd", key+".pass")
	env.RegisterAlias(key+".addr", key+".host")
	return
}

// Validate checks that the configuration is valid.
func (c Config) Validate() (err error) {
	if c.Addr == "" {
		err = errors.Append(err, errors.New("database addr is required"))
	}

	if c.User == "" {
		err = errors.Append(err, errors.New("database user is required"))
	}

	if c.DBName == "" {
		err = errors.Append(err, errors.New("database name is required"))
	}

	return
}
