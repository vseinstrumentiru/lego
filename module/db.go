package module

import (
	"github.com/vseinstrumentiru/lego/v2/di"
	"github.com/vseinstrumentiru/lego/v2/transport/mongodb"
	"github.com/vseinstrumentiru/lego/v2/transport/mysql"
	"github.com/vseinstrumentiru/lego/v2/transport/postgres"
	"github.com/vseinstrumentiru/lego/v2/transport/sql"
)

func MySQLPack() []di.Module {
	return []di.Module{
		MySQLConnector,
		SQLDB,
		HealthChecker,
	}
}

func PostgresPack() []di.Module {
	return []di.Module{
		PostgresConnector,
		SQLDB,
		HealthChecker,
	}
}

//nolint:stylecheck
// Deprecated: use MongoDBPack.
func MongoDbPack() []di.Module {
	return MongoDBPack()
}

func MongoDBPack() []di.Module {
	return []di.Module{
		MongoDBConnector,
		HealthChecker,
	}
}

func SQLDB() (interface{}, []interface{}) {
	return sql.Provide, nil
}

func MySQLConnector() (interface{}, []interface{}) {
	return mysql.Provide, nil
}

func PostgresConnector() (interface{}, []interface{}) {
	return postgres.Provide, nil
}

//nolint:stylecheck
// Deprecated: use MongoDBConnector.
func MongoDbConnector() (interface{}, []interface{}) {
	return MongoDBConnector()
}

func MongoDBConnector() (interface{}, []interface{}) {
	return mongodb.Provide, nil
}
