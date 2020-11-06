package database

import (
	"database/sql"

	"emperror.dev/emperror"

	"github.com/vseinstrumentiru/lego/v2/internal/deprecated"
	"github.com/vseinstrumentiru/lego/v2/server/shutdown"
)

// Deprecated: use DI.
func NewSQLConnection(_ interface{}, _ interface{}) (conn *sql.DB, closer *shutdown.CloseGroup) {
	err := deprecated.Container.Execute(func(i *sql.DB, j *shutdown.CloseGroup) {
		conn = i
		closer = j
	})

	emperror.Panic(err)

	return
}
