package sql

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/errors"
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"go.uber.org/dig"

	"github.com/vseinstrumentiru/lego/server/shutdown"
	"github.com/vseinstrumentiru/lego/transport/database"
)

type args struct {
	dig.In
	MySQL *database.MySQLConnector `optional:"true"`

	Closer *shutdown.CloseGroup
	Health health.Health
}

func Provide(in args) (*sql.DB, error) {
	var connector driver.Connector

	if in.MySQL != nil {
		connector = in.MySQL
	} else {
		return nil, errors.New("connector not found. you must provide connector")
	}

	conn := sql.OpenDB(connector)
	stopStats := ocsql.RecordStats(conn, 5*time.Second)

	err := in.Health.RegisterCheck(&health.Config{
		Check:           checks.Must(checks.NewPingCheck("db.check", conn, time.Millisecond*100)),
		ExecutionPeriod: 3 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	in.Closer.Add(shutdown.SimpleCloseFn(stopStats))
	in.Closer.Add(conn)

	return conn, nil
}
