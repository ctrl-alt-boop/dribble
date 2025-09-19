package result

import (
	"database/sql"
)

type Kind int

const (
	KindNone Kind = iota
	KindRow
	KindScalar
	KindList
	KindTable
	KindSet
)

type (
	Scalar struct {
		Value   any
		DBField string
	}

	Object struct {
		Fields map[string]any
		Type   string
	}
)

func RowToScalar(rows *sql.Rows) Scalar {
	return Scalar{}
}

func RowsToObject(rows *sql.Rows) (Object, error) {
	return Object{}, nil
}
