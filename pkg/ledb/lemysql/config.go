package lemysql

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Host string
	Port int
	User string
	Pass string
	Name string

	Params map[string]string
}

func (c Config) SetDefaults(env *viper.Viper, flag *pflag.FlagSet) {
	return
}

// Validate checks that the configuration is valid.
func (c Config) Validate() (err error) {
	if c.Host == "" {
		err = errors.Append(err, errors.New("database host is required"))
	}

	if c.Port == 0 {
		err = errors.Append(err, errors.New("database port is required"))
	}

	if c.User == "" {
		err = errors.Append(err, errors.New("database user is required"))
	}

	if c.Name == "" {
		err = errors.Append(err, errors.New("database name is required"))
	}

	return
}

// DSN returns a MySQL driver compatible data source name.
func (c Config) DSN() string {
	var params string

	if len(c.Params) > 0 {
		var query string

		for key, value := range c.Params {
			if query != "" {
				query += "&"
			}

			query += key + "=" + value
		}

		params = "?" + query
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
		params,
	)
}
