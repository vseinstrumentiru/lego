package mysql

import (
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type Args struct {
	dig.In
	Config *Config
	Logger multilog.Logger
}

func ProvideConnector(in Args) (*Connector, error) {
	conn, err := Provide(in)

	if err != nil {
		return nil, err
	}

	return &Connector{conn}, nil
}

func Provide(in Args) (driver.Connector, error) {
	connector, err := mysql.NewConnector(&in.Config.Config)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	logger := in.Logger.WithFields(map[string]interface{}{"component": "mysql"})
	_ = mysql.SetLogger(logur.NewErrorPrintLogger(logger))

	return ocsql.WrapConnector(
		connector,
		ocsql.WithOptions(in.Config.Trace),
	), nil
}
