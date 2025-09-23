package result

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type List struct {
	Values     []any
	DBTypeName string
	DBField    string
	ScanType   reflect.Type
	ActualType reflect.Type
	DBKind     string
}

// String implements Response.
func (l List) String() string {
	stringValues := make([]string, len(l.Values))
	for i, v := range l.Values {
		if v == nil {
			stringValues[i] = "NULL"
		} else {
			stringValues[i] = fmt.Sprintf("%v", v)
		}
	}
	return strings.Join(stringValues, "\n")
}

// Value implements Response.
func (l List) Get() any {
	return l.get()
}

func (l List) get() []any {
	return l.Values
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
