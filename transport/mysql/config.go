package mysql

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"

	"github.com/vseinstrumentiru/lego/config"
)

type Config struct {
	mysql.Config `mapstructure:",squash"`
	Trace        ocsql.TraceOptions
}

func (c Config) SetDefaults(env config.Env) {
	env.SetAlias("pass", "passwd")

	env.SetDefault("parseTime", true)
	env.SetDefault("rejectReadOnly", true)
	env.SetDefault("allowNativePasswords", true)
	env.SetDefault("trace.query", true)
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
