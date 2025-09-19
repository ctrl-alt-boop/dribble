package result

import (
	"database/sql"
	"reflect"
)

type List struct {
	Values     []any
	DBTypeName string
	DBField    string
	ScanType   reflect.Type
	ActualType reflect.Type
	DBKind     string
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
