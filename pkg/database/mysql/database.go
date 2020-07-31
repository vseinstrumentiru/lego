package mysql

import (
	"database/sql/driver"
	"github.com/vseinstrumentiru/lego/tools/databasetools/mysqlconnector"
)

// NewConnector returns a new database connector for the application.
// deprecated
func NewConnector(config Config) (driver.Connector, error) {
	return mysqlconnector.NewConnector(config)
}
