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

func GetSupportedDrivers() []string {
	return supportedDrivers
}

func GetDriverDefaults() map[string]*connection.Settings {
	return map[string]*connection.Settings{
		database.DriverPostgreSQL: &postgresDefault,
		database.DriverMySql:      &mysqlDefault,
		database.DriverSQLite:     &sqliteDefault,
	}
}

var postgresDefault = connection.Settings{
	Type:       connection.Driver,
	DriverName: "postgres",
	Ip:         "127.0.0.1",
	Port:       5432,
	AdditionalSettings: map[string]string{
		"sslmode": "disable",
	},
}

var mysqlDefault = connection.Settings{
	Type:       connection.Driver,
	DriverName: "mysql",
	Ip:         "127.0.0.1",
	Port:       3306,
}

var sqliteDefault = connection.Settings{
	Type:       connection.Driver,
	DriverName: "sqlite",
}
