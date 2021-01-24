package sql

import (
	"context"
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
)

func NewConnector(drv driver.Driver, dsn string, opts *ocsql.TraceOptions) driver.Connector {
	return &connector{
		dsn:     dsn,
		driver:  drv,
		options: opts,
	}
}

type connector struct {
	dsn     string
	driver  driver.Driver
	options *ocsql.TraceOptions
}

func (c *connector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.dsn)
}

func (c *connector) Driver() driver.Driver {
	if c.options == nil {
		return c.driver
	}

	return ocsql.Wrap(
		c.driver,
		ocsql.WithOptions(*c.options),
	)
}
