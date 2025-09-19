package database

import (
	"fmt"
	"slices"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
)

func CreateClientForDialect(dialect target.Dialect) (database.Database, error) {
	switch {
	case slices.Contains(sql.SupportedDialects, dialect):
		return sql.NewClient(dialect), nil
	case slices.Contains(nosql.SupportedDialects, dialect):
		return nosql.NewClient(dialect), nil
	default:
		return nil, fmt.Errorf("unknown or unsupported driver: %s", dialect)
	}
}
