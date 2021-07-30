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
)

type Args struct {
	dig.In

	Connector driver.Connector
	Closer    *shutdown.CloseGroup `optional:"true"`
	Health    health.Health        `optional:"true"`
}

func Provide(in Args) (*sql.DB, error) {
	if in.Connector == nil {
		return nil, errors.New("connector not found. you must provide `driver.Connector`")
	}

	conn := sql.OpenDB(in.Connector)
	stopStats := ocsql.RecordStats(conn, 5*time.Second)

	if in.Health != nil {
		err := in.Health.RegisterCheck(&health.Config{
			Check:           checks.Must(checks.NewPingCheck("db.check", conn, time.Millisecond*100)),
			ExecutionPeriod: 3 * time.Second,
		})
		if err != nil {
			return nil, err
		}

		if in.Closer != nil {
			in.Closer.Add(shutdown.SimpleCloseFn(stopStats))
		}
	}

	if in.Closer != nil {
		in.Closer.Add(conn)
	}

	return conn, nil
}
