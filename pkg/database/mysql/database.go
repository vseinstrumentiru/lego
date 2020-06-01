package mysql

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	"database/sql/driver"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"
)

// NewConnector returns a new database connector for the application.
func NewConnector(config Config) (driver.Connector, error) {
	connector, err := mysql.NewConnector(&config.Config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ocsql.WrapConnector(
		connector,
		ocsql.WithOptions(config.Trace),
	), nil
}
