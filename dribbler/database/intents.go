package database

import (
	"github.com/ctrl-alt-boop/dribble/database"
)

func NewDatabaseNamesQuery() *database.Intent {
	builder := database.Select("db_name").From("information_schema.schemata")
	return builder.ToQuery()
}

func ExecuteDatabaseNamesQuery() []string {

	return nil
}
