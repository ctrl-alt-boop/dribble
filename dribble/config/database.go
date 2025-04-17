package config

import "github.com/ctrl-alt-boop/gooldb/pkg/connection"

var SavedConfigs map[string]*connection.Settings = map[string]*connection.Settings{
	"postgres_win": connection.NewSettings(
		connection.WithDriver("postgres"),
		connection.WithHost("172.24.208.1", 5432),
		connection.WithUser("valmatics"),
		connection.WithPassword("valmatics"),
		connection.WithSetting("sslmode", "disable"),
	),
	"postgres_local": connection.NewSettings(
		connection.WithDriver("postgres"),
		connection.WithHost("localhost", 5432),
		connection.WithUser("postgres_user"),
		connection.WithPassword("postgres_user"),
		connection.WithSetting("sslmode", "disable"),
	),
	"mysql_local": connection.NewSettings(
		connection.WithDriver("mysql"),
		connection.WithHost("localhost", 3306),
		connection.WithUser("mysql_user"),
		connection.WithPassword("mysql_user"),
	),
}
