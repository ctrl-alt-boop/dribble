package config

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

type Connection struct {
	Type               connection.Type
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

var SavedConfigs map[string]*connection.Settings = map[string]*connection.Settings{
	"postgres_win": connection.NewSettings(
		connection.AsType(connection.Server),
		connection.WithDriver("postgres"),
		connection.WithHost("172.24.208.1", 5432),
		connection.WithUser("valmatics"),
		connection.WithPassword("valmatics"),
		connection.WithSetting("sslmode", "disable"),
	),
	"postgres_local": connection.NewSettings(
		connection.AsType(connection.Server),
		connection.WithDriver("postgres"),
		connection.WithHost("localhost", 5432),
		connection.WithUser("postgres_user"),
		connection.WithPassword("postgres_user"),
		connection.WithSetting("sslmode", "disable"),
	),
	"mysql_local": connection.NewSettings(
		connection.AsType(connection.Server),
		connection.WithDriver("mysql"),
		connection.WithHost("localhost", 3306),
		connection.WithUser("mysql_user"),
		connection.WithPassword("mysql_user"),
	),
}

func GetDriverDefaults() map[string]Connection {
	defaults := gooldb.GetDriverDefaults()
	driverDefaults := make(map[string]Connection)
	for name, settings := range defaults {
		driverDefaults[name] = Connection{
			Type:               connection.Driver,
			DriverName:         name,
			Ip:                 settings.Ip,
			Port:               settings.Port,
			Username:           settings.Username,
			Password:           settings.Password,
			AdditionalSettings: settings.AdditionalSettings,
		}
	}
	return driverDefaults
}
