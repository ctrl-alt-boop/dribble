package result

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ctrl-alt-boop/dribble/datasource"
)

type (
	Column struct {
		Name     string
		ScanType reflect.Type
		DBType   string
	}

	Row struct {
		Values []any
	}

	Table struct {
		columns []*Column
		rows    []*Row

		Resolver datasource.SQLAdapter
	}
)

func (r Row) String() string {
	return strings.Join(sliceTransform(r.Values, func(value any) string {
		return fmt.Sprint(value)
	}), ", ")
}

// Value implements Response.
func (r Row) Get() any {
	return r.get()
}

func (r Row) get() []any {
	return r.Values
}

func (dt Table) String() string {
	tableString := ""
	for _, row := range dt.GetRowStringsAll() {
		fmt.Println(row)
		tableString += strings.Join(row, ", ")
	}
	return tableString
}

// Value implements Response.
func (dt Table) Get() any {
	return dt.get()
}

func (dt Table) get() Table {
	return dt
}

func (dt *Table) NumColumns() int {
	return len(dt.columns)
}

func (dt *Table) NumRows() int {
	return len(dt.rows)
}

func (dt *Table) Columns() []*Column {
	return dt.columns
}

func (dt *Table) Rows() []*Row {
	return dt.rows
}

func CreateDataScalar(rows []Row) any {
	return rows[0].Values[0]
}

func CreateDataList(rows []Row) []any {
	return sliceTransform(rows, func(row Row) any {
		return row.Values[0]
	})
}

func NewTable(columns []*Column, rows []*Row) *Table {
	return &Table{
		columns: columns,
		rows:    rows,
	}
}

func (dt *Table) GetRowColumn(row, column int) (string, error) { // Needs a lookielook to see if other drivers are at least similar to this
	rowColumn := dt.rows[row].Values[column]
	switch value := rowColumn.(type) {
	case string, int, int32, int64, float32, float64, uint, bool:
		return fmt.Sprint(value), nil
	case time.Time:
		return fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07")), nil
	case []byte:
		resolved, err := dt.Resolver.ResolveType(dt.columns[column].DBType, value)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", resolved), nil
		}
	case nil:
		return "null", nil
	default:
		err := fmt.Errorf("unknown value type %T for %s", value, dt.columns[column].DBType)
		return "", err
	}
}

func (dt *Table) GetRowStringsAll() [][]string {
	rows := make([][]string, len(dt.columns))
	for i := range dt.NumRows() {
		rows[i] = dt.GetRowStrings(i)
	}

	return rows
}

func (dt *Table) GetRowStrings(index int) []string {
	row := make([]string, len(dt.columns))
	for columnIndex := range dt.columns {
		value, err := dt.GetRowColumn(index, columnIndex)
		if err != nil {
			row[columnIndex] = err.Error()
		} else {
			row[columnIndex] = value
		}
	}
	return row
}

func (dt *Table) GetRowString(index int) string {
	row := dt.GetRowStrings(index)
	return strings.Join(row, " | ")
}

func (dt *Table) GetColumnRows(columnIndex int) (rows []string, columnWidth int) {
	columnRows := make([]string, dt.NumRows())
	for rowIndex := range dt.rows {
		value, err := dt.GetRowColumn(rowIndex, columnIndex)
		if err != nil {
			columnRows[rowIndex] = err.Error()
		} else {
			columnRows[rowIndex] = value
		}

		columnWidth = max(columnWidth, len(columnRows[rowIndex]))
	}
	return columnRows, columnWidth
}

func (dt *Table) ColumnSlices() (names []string, types []string, dbTypes []string) {
	names = dt.ColumnNames()
	types = dt.ColumnTypeStrings()
	dbTypes = dt.ColumnDatabaseTypeStrings()
	return
}

func (dt *Table) ColumnNames() []string {
	return sliceTransform(dt.columns, func(col *Column) string {
		return col.Name
	})
}

func (dt *Table) ColumnTypeStrings() []string {
	return sliceTransform(dt.columns, func(col *Column) string {
		return col.ScanType.Kind().String()
	})
}

func (dt *Table) ColumnDatabaseTypeStrings() []string {
	return sliceTransform(dt.columns, func(col *Column) string {
		return col.DBType
	})
}

func (dt *Table) ClearRows() error {
	dt.rows = make([]*Row, 0)
	return nil
}

func sliceTransform[T any, U any](slice []T, selector func(T) U) []U {
	results := make([]U, len(slice))
	for i, value := range slice {
		results[i] = selector(value)
	}
	return results
}
