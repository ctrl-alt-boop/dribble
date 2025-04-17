package gooldb

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

var supportedDrivers []string = []string{
	database.DriverPostgreSQL,
	database.DriverMySql,
	database.DriverSQLite,
}

func (gool *GoolDb) GetSupportedDrivers() []string {
	return supportedDrivers
}

func (gool *GoolDb) GetDriverDefaults() map[string]*connection.Settings {
	return map[string]*connection.Settings{
		database.DriverPostgreSQL: &postgresDefault,
		database.DriverMySql:      &mysqlDefault,
		database.DriverSQLite:     &sqliteDefault,
	}
}

var postgresDefault = connection.Settings{
	DriverName: "postgres",
	Ip:         "127.0.0.1",
	Port:       5432,
	AdditionalSettings: map[string]string{
		"sslmode": "disable",
	},
}

var mysqlDefault = connection.Settings{
	DriverName: "mysql",
	Ip:         "127.0.0.1",
	Port:       3306,
}

var sqliteDefault = connection.Settings{
	DriverName: "sqlite",
}
