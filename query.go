package dribble

import (
	"github.com/ctrl-alt-boop/dribble/database"
)

type (
	Query interface {
		ToSQL(dialect database.Dialect) (queryString string, params []any, err error)
		ToSQLFormatted(dialect database.Dialect) (queryString string, params []any, err error)
		Parameters() []any
	}

	Dependency struct {
		SourceId    int
		SouceDataId int

		TargetId     int
		TargetDataId int
	}

	Batch struct {
		Queries      []Query
		Dependencies []Dependency
	}
)
