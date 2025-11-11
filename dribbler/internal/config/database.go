package config

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/dsn"
	"github.com/ctrl-alt-boop/dribble/target"
)

type Connection struct {
	Type               target.Type
	DriverName         string
	Ip                 string
	Port               int
	Username           string
	Password           string
	AdditionalSettings map[string]string
}

func (c Connection) String() string {
	return fmt.Sprintf("%s://%s:%d", c.DriverName, c.Ip, c.Port)
}

var SavedConfigs map[string]datasource.Namer = map[string]datasource.Namer{
	"postgres_win": dsn.PostgresDSN(
		dsn.PostgresAddr("172.24.208.1"),
		dsn.PostgresPort(5432),
		dsn.PostgresUsername("valmatics"),
		dsn.PostgresPassword("valmatics"),
		dsn.PostgresSSLMode(dsn.SSLModeDisable),
	),
	"postgres_local": dsn.PostgresDSN(
		dsn.PostgresAddr("localhost"),
		dsn.PostgresPort(5432),
		dsn.PostgresUsername("postgres_user"),
		dsn.PostgresPassword("postgres_user"),
		dsn.PostgresSSLMode(dsn.SSLModeDisable),
	),
	"mysql_local": dsn.MySQLDSN(
		dsn.MySQLAddr("localhost"),
		dsn.MySQLPort(3306),
		dsn.MySQLUsername("mysql_user"),
		dsn.MySQLPassword("mysql_user"),
	),
}

func GetDriverDefaults() map[string]Connection {
	// defaults := dribble.GetDriverDefaults()
	driverDefaults := make(map[string]Connection)
	// for name, settings := range defaults {
	// 	driverDefaults[name] = Connection{
	// 		Type:               target.TypeDriver,
	// 		DriverName:         name,
	// 		Ip:                 settings.Ip,
	// 		Port:               settings.Port,
	// 		Username:           settings.Username,
	// 		Password:           settings.Password,
	// 		AdditionalSettings: settings.AdditionalSettings,
	// 	}
	// }
	return driverDefaults
}
