package target

import (
	"maps"

	"github.com/ctrl-alt-boop/dribble/database"
)

type Option func(*Target)

func WithConnectionProperties(properties *database.ConnectionProperties) Option {
	return func(target *Target) {
		target.Properties.Addr = properties.Addr
		target.Properties.Port = properties.Port
		target.Properties.DBName = properties.DBName
		target.Properties.Username = properties.Username
		target.Properties.Password = properties.Password
		maps.Copy(target.Properties.Additional, properties.Additional)
	}
}

func AsTableSelect(table string) Option {
	return func(target *Target) {
		target.Type = TypeDatabase
		prop := target.Properties
		prop.Additional["select"] = "SELECT * FROM " + table
	}
}

func AsQuery(query string) Option {
	return func(target *Target) {
		target.Type = TypeDatabase
		target.Properties.Additional["query"] = query
	}
}
