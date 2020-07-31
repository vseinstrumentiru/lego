package databasetools

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	"database/sql"
	"database/sql/driver"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/facebookincubator/ent/dialect"
	entsql "github.com/facebookincubator/ent/dialect/sql"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/tools/databasetools/mysqlconnector"
	"reflect"
	"time"
)

func NewSQLConnection(p lego.Process, config interface{}) (conn *sql.DB, closer lego.CloserGroup) {
	var connector driver.Connector
	var err error

	switch cfg := config.(type) {
	case mysqlconnector.Config:
		connector, err = mysqlconnector.NewConnector(cfg)
		emperror.Panic(err)
	default:
		emperror.Panic(errors.NewWithDetails("undefined sql database config type", "type", reflect.TypeOf(config)))
	}

	conn = sql.OpenDB(connector)
	closer.Add(conn)
	stopStats := ocsql.RecordStats(conn, 5*time.Second)
	// Record DB stats every 5 seconds until we exit
	closer.Add(lego.CloseFn(func() error {
		stopStats()
		return nil
	}))

	_ = p.RegisterCheck(&health.Config{
		Check:           checks.Must(checks.NewPingCheck("db.check", conn, time.Millisecond*100)),
		ExecutionPeriod: 3 * time.Second,
	})

	return
}

func NewEntDriver(p lego.Process, config interface{}) (drv dialect.Driver, closer lego.CloserGroup) {
	var conn *sql.DB

	conn, closer = NewSQLConnection(p, config)
	var dialect string

	switch config.(type) {
	case mysqlconnector.Config:
		dialect = "mysql"
	default:
		emperror.Panic(errors.NewWithDetails("undefined sql dialect type", "type", reflect.TypeOf(config)))
	}

	drv = entsql.OpenDB(dialect, conn)
	closer.Add(drv)

	return drv, closer
}
