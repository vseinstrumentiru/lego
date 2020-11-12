package mysql

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/multilog"
)

type args struct {
	dig.In
	Config *Config
	Logger multilog.Logger
}

func Provide(in args) (*Connector, error) {
	connector, err := mysql.NewConnector(&in.Config.Config)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	logger := in.Logger.WithFields(map[string]interface{}{"component": "mysql"})
	_ = mysql.SetLogger(logur.NewErrorPrintLogger(logger))

	return &Connector{
		Connector: ocsql.WrapConnector(
			connector,
			ocsql.WithOptions(in.Config.Trace),
		),
	}, nil
}
