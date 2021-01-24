package mysql

import (
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"

	"github.com/vseinstrumentiru/lego/v2/config"
)

type Config struct {
	mysql.Config `mapstructure:",squash"`
}

func (c Config) SetDefaults(env config.Env) {
	env.SetAlias("pass", "passwd")

	env.SetDefault("parseTime", true)
	env.SetDefault("rejectReadOnly", true)
	env.SetDefault("allowNativePasswords", true)
}

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
