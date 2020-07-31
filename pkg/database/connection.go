package database

import (
	"database/sql"
	"github.com/facebookincubator/ent/dialect"
	"github.com/vseinstrumentiru/lego/internal/lego"
	"github.com/vseinstrumentiru/lego/tools/databasetools"
)

// deprecated
func NewSQLConnection(p lego.Process, config interface{}) (conn *sql.DB, closer lego.CloserGroup) {
	return databasetools.NewSQLConnection(p, config)
}

// deprecated
func NewEntDriver(p lego.Process, config interface{}) (drv dialect.Driver, closer lego.CloserGroup) {
	return databasetools.NewEntDriver(p, config)
}
