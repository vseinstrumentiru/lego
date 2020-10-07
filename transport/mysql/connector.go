package mysql

import "database/sql/driver"

type Connector interface {
	driver.Connector
}
