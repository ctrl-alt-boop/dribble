package config

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
)

type Connection struct {
	Type               database.TargetType
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

var SavedConfigs map[string]*database.Target = map[string]*database.Target{
	"postgres_win": database.NewTarget("postgres_win", database.TargetServer,
		database.WithDriver("postgres"),
		database.WithHost("172.24.208.1", 5432),
		database.WithUser("valmatics"),
		database.WithPassword("valmatics"),
		database.WithSetting("sslmode", "disable"),
	),
	"postgres_local": database.NewTarget("postgres_local", database.TargetServer,
		database.WithDriver("postgres"),
		database.WithHost("localhost", 5432),
		database.WithUser("postgres_user"),
		database.WithPassword("postgres_user"),
		database.WithSetting("sslmode", "disable"),
	),
	"mysql_local": database.NewTarget("mysql_local", database.TargetServer,
		database.WithDriver("mysql"),
		database.WithHost("localhost", 3306),
		database.WithUser("mysql_user"),
		database.WithPassword("mysql_user"),
	),
}

func GetDriverDefaults() map[string]Connection {
	defaults := dribble.GetDriverDefaults()
	driverDefaults := make(map[string]Connection)
	for name, settings := range defaults {
		driverDefaults[name] = Connection{
			Type:               database.TargetDriver,
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
