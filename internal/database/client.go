package database

import (
	"errors"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/nosql"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql"
)

func CreateClientForType(t database.Type) (database.Database, error) {
	switch t := t.(type) {
	case database.SQLDialectType:
		return CreateSQLClient(t)
	case database.NoSQLModelType:
		return CreateNoSQLClient(t)
	case database.GraphType:
		return nil, nil
	case database.TimeSeriesType:
		return nil, nil
	default:
		return nil, errors.New("unknown type")
	}
}

func CreateSQLClient(dialect database.SQLDialectType) (database.SQL, error) {
	executor, err := sql.NewExecutor(dialect)
	if err != nil {
		return nil, err
	}
	return executor, nil
}

func CreateNoSQLClient(modelType database.NoSQLModelType) (database.NoSQL, error) {
	executor, err := nosql.NewExecutor(modelType)
	if err != nil {
		return nil, err
	}
	return executor, nil
}
