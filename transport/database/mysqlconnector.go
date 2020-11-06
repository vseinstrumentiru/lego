package database

import "database/sql/driver"

type MySQLConnector struct {
	driver.Connector
}
