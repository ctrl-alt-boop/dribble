package result

import (
	"database/sql"
	"reflect"

	"github.com/ctrl-alt-boop/dribble/database"
)

type (
	List struct {
		Values     []any
		DBTypeName string
		DBField    string
		ScanType   reflect.Type
		ActualType reflect.Type
		DBKind     string
	}

	Scalar struct {
		Value   any
		DBField string
	}

	Object struct {
		Fields map[string]any
		Type   string
	}

	Kind int
)

const (
	KindNone Kind = iota - 1
	KindRow
	KindScalar
	KindList
	KindTable
	KindSet
)

var DefaultOperationResults = map[database.OperationType]Kind{
	database.Read:    KindTable,
	database.Create:  KindScalar,
	database.Update:  KindScalar,
	database.Delete:  KindScalar,
	database.Execute: KindScalar,
}

// func RowsToList(rows *sql.Rows) (List, error) {
// 	return List{}, nil
// }

func RowToScalar(rows *sql.Rows) Scalar {
	return Scalar{}
}

func RowsToObject(rows *sql.Rows) (Object, error) {
	return Object{}, nil
}

func RowsToList(dbRows *sql.Rows) List {
	dbColumns, err := dbRows.ColumnTypes()
	if err != nil {
		return List{}
	}

	if len(dbColumns) > 1 {
		return List{}
	}
	column := dbColumns[0]
	list := List{
		Values:     make([]any, 0),
		DBField:    column.Name(),
		ScanType:   column.ScanType(),
		DBTypeName: column.DatabaseTypeName(),
		DBKind:     column.ScanType().Kind().String(),
	}

	for dbRows.Next() {
		rowValuePtr := reflect.New(list.ScanType).Interface()
		err := dbRows.Scan(rowValuePtr)
		if err != nil {
			continue
		}
		scannedValue := reflect.ValueOf(rowValuePtr).Elem().Interface()
		list.ActualType = reflect.TypeOf(scannedValue)
		switch scannedValue := scannedValue.(type) {
		case []byte:
			list.Values = append(list.Values, string(scannedValue))
		default:
			list.Values = append(list.Values, scannedValue)
		}
	}
	return list
}
