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

func MongoDbPack() []di.Module {
	return []di.Module{
		MongoDbConnector,
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

func MongoDbConnector() (interface{}, []interface{}) {
	return mongodb.Provide, nil
}
