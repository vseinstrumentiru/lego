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

	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
	"github.com/vseinstrumentiru/lego/v2/transport/mysql"
)

type Args struct {
	dig.In
	Connector driver.Connector `optional:"true"`
	// Deprecated: use driver.Connector
	MySQL *mysql.Connector `optional:"true"`

	Closer *shutdown.CloseGroup
	Health health.Health
}

func Provide(in Args) (*sql.DB, error) {
	if in.Connector == nil {
		if in.MySQL != nil {
			in.Connector = in.MySQL.Connector
		} else {
			return nil, errors.New("connector not found. you must provide `driver.Connector`")
		}
	}

	conn := sql.OpenDB(in.Connector)
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
