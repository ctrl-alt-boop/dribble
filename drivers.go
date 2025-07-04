package dribble

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/nosql"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql"
)

type DriverName = string

const (
	DriverPostgreSQL DriverName = "postgres"
	DriverMySql      DriverName = "mysql"
	DriverSQLite     DriverName = "sqlite3"
	DriverMongoDB    DriverName = "mongodb"
)

var supportedDrivers []string = []string{
	DriverPostgreSQL,
	DriverMySql,
	DriverSQLite,
}

var (
	_ database.Driver = &sql.MySql{}
	_ database.Driver = &sql.Postgres{}
	_ database.Driver = &sql.SQLite3{}
)

func CreateDriverFromTarget(target *database.Target) (database.Driver, error) {
	switch DriverName(target.DriverName) {
	case DriverMySql:
		return sql.NewMySqlDriver(target)
	case DriverPostgreSQL:
		return sql.NewPostgresDriver(target)
	case DriverSQLite:
		return sql.NewSQLite3Driver(target)
	case DriverMongoDB:
		return nosql.NewMongoDBDriver(target)
	default:
		return nil, fmt.Errorf("unknown or unsupported driver: %s", target.DriverName)
	}
}

func GetSupportedDrivers() []string {
	return supportedDrivers
}

func GetDriverDefaults() map[string]*database.Target {
	return map[string]*database.Target{
		DriverPostgreSQL: &postgresDefault,
		DriverMySql:      &mysqlDefault,
		DriverSQLite:     &sqliteDefault,
	}
}

var postgresDefault = database.Target{
	Type:       database.DBDriver,
	DriverName: "postgres",
	Ip:         "127.0.0.1",
	Port:       5432,
	AdditionalSettings: map[string]string{
		"sslmode": "disable",
	},
}

var mysqlDefault = database.Target{
	Type:       database.DBDriver,
	DriverName: "mysql",
	Ip:         "127.0.0.1",
	Port:       3306,
}

var sqliteDefault = database.Target{
	Type:       database.DBDriver,
	DriverName: "sqlite",
}
