package result

import (
	"database/sql"
	"fmt"
)

type Kind int

const (
	KindNone Kind = iota
	KindRow
	KindScalar
	KindList
	KindTable
	KindSet
	KindObject
	KindMap
)

var (
	_ Body = Row{}
	_ Body = List{}
	_ Body = Table{}
)

type (
	Scalar struct {
		value   int
		DBField string
	}

	Object struct {
		Fields map[string]any
		Type   string
	}

	Body interface {
		fmt.Stringer
		Get() any
	}
)

func RowToScalar(row *sql.Row) (*Scalar, error) {
	var scalar int
	err := row.Scan(&scalar)
	if err != nil {
		return &Scalar{value: 0}, fmt.Errorf("error scanning row: %w", err)
	}
	return &Scalar{value: scalar}, nil
}

func RowsToObject(rows *sql.Rows) (Object, error) {
	return Object{}, nil
}

// Please don't use
func RowToRow(sqlRow *sql.Row) (Row, error) {
	row := Row{
		Values: nil,
	}
	scanArr := make([]any, 25) // I don't know yet how to improve this
	for i := range row.Values {
		scanArr[i] = &row.Values[i]
	}

	err := sqlRow.Scan(scanArr...)
	if err != nil {
		return row, fmt.Errorf("error scanning row: %w", err)
	}
	return row, nil
}

func (s Scalar) String() string {
	return fmt.Sprintf("%d", s.value)
}

func (s Scalar) Get() any {
	return s.get()
}

func (s Scalar) get() int {
	return s.value
}

func (o Object) String() string {
	return fmt.Sprintf("%v", o.Fields)
}

func (o Object) Get() any {
	return o.Fields
}
