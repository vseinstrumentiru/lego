package mysql

import "database/sql/driver"

type Connector struct {
	driver.Connector
}
