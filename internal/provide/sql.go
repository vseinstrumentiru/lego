package provide

import (
	"github.com/vseinstrumentiru/lego/v2/transport/mysql"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"
)

func Sql() []interface{} {
	return []interface{}{
		mysql.ProvideConnector,
		sql.Provide,
	}
}
