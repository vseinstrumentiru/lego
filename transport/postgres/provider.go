package postgres

import (
	"database/sql/driver"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/newrelic/sqlparse"
	"strconv"

	"contrib.go.opencensus.io/integrations/ocsql"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/v2/metrics/tracing"
	"github.com/vseinstrumentiru/lego/v2/multilog"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type Args struct {
	dig.In
	Config   *Config
	Trace    *tracing.Config `optional:"true"`
	Newrelic *newrelic.Application
	Logger   multilog.Logger `optional:"true"`
}

func Provide(in Args) (driver.Connector, error) {
	config, err := pgx.ParseConfig(in.Config.DSN)
	if err != nil {
		return nil, err
	}

	if in.Logger != nil {
		config.Logger = &logger{
			Logger: in.Logger.WithFields(map[string]interface{}{"component": "postgresql"}),
		}
	}

	dsn := stdlib.RegisterConnConfig(config)

	var options *ocsql.TraceOptions
	if in.Trace != nil && in.Trace.SQL != nil {
		options = in.Trace.SQL
	}

	connector := sql.NewConnector(stdlib.GetDefaultDriver(), dsn, options)

	if in.Newrelic != nil {
		builder := newrelic.SQLDriverSegmentBuilder{
			BaseSegment: newrelic.DatastoreSegment{
				Product:      newrelic.DatastorePostgres,
				DatabaseName: config.Database,
				Host:         config.Host,
				PortPathOrID: strconv.Itoa(int(config.Port)),
			},
			ParseQuery: sqlparse.ParseQuery,
		}

		connector = newrelic.InstrumentSQLConnector(sql.NewConnector(stdlib.GetDefaultDriver(), dsn, options), builder)
	}

	return connector, nil
}
