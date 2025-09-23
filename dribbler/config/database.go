package config

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/dsn/mysql"
	"github.com/ctrl-alt-boop/dribble/dsn/postgres"
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

var SavedConfigs map[string]database.DataSourceNamer = map[string]database.DataSourceNamer{
	"postgres_win": postgres.NewDSN(
		postgres.WithAddr("172.24.208.1"),
		postgres.WithPort(5432),
		postgres.WithUsername("valmatics"),
		postgres.WithPassword("valmatics"),
		postgres.WithSSLMode(postgres.SSLModeDisable),
	),
	"postgres_local": postgres.NewDSN(
		postgres.WithAddr("localhost"),
		postgres.WithPort(5432),
		postgres.WithUsername("postgres_user"),
		postgres.WithPassword("postgres_user"),
		postgres.WithSSLMode(postgres.SSLModeDisable),
	),
	"mysql_local": mysql.NewDSN(
		mysql.WithAddr("localhost"),
		mysql.WithPort(3306),
		mysql.WithUsername("mysql_user"),
		mysql.WithPassword("mysql_user"),
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
