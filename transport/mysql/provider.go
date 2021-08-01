package mysql

import (
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/dig"
	"logur.dev/logur"

	"github.com/vseinstrumentiru/lego/v2/log"
	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"
)

type Args struct {
	dig.In
	Config *Config
	Trace  *tracing.Config `optional:"true"`
	Logger log.Logger
}

func ProvideConnector(in Args) (*Connector, error) {
	conn, err := Provide(in)
	if err != nil {
		return nil, err
	}

	return &Connector{conn}, nil
}

func Provide(in Args) (driver.Connector, error) {
	logger := in.Logger.WithFields(map[string]interface{}{"component": "mysql"})
	_ = mysql.SetLogger(logur.NewErrorPrintLogger(logger))

	var options *ocsql.TraceOptions
	if in.Trace != nil && in.Trace.SQL != nil {
		options = in.Trace.SQL
	}

	dsn := in.Config.FormatDSN()

	return sql.NewConnector(mysql.MySQLDriver{}, dsn, options), nil
}
