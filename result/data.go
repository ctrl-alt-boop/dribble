package result

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ctrl-alt-boop/dribble/database"
)

func ParseRows(dbRows *sql.Rows) ([]Column, []Row) {
	dbColumns, err := dbRows.ColumnTypes()
	if err != nil {
		return nil, nil
	}
	columns := make([]Column, len(dbColumns))
	for i := range dbColumns {
		columns[i] = Column{
			Name:     dbColumns[i].Name(),
			ScanType: dbColumns[i].ScanType(),
			DBType:   dbColumns[i].DatabaseTypeName(),
		}
	}

	rows := make([]Row, 0)
	for dbRows.Next() {
		row := Row{
			Values: make([]any, len(dbColumns)),
		}
		scanArr := make([]any, len(dbColumns))
		for i := range row.Values {
			scanArr[i] = &row.Values[i]
		}

		err := dbRows.Scan(scanArr...)
		if err != nil {
			continue
		}
		// for i := range row.Values {
		// 	row.Values[i], err = ResolveTypes(resolver, row.Values[i], columns[i])
		// 	if err != nil {
		// 		continue
		// 	}
		// }

		rows = append(rows, row)
	}
	return columns, rows
}

func ResolveTypes(resolver database.Dialect, rowValue any, column Column) (any, error) {
	switch value := rowValue.(type) {
	case string, int, int32, int64, float32, float64, uint, bool:
		return value, nil
	case time.Time:
		return fmt.Sprint(value.Format("2006-01-02 15:04:05.000000-07")), nil
	case []byte:
		resolved, err := resolver.ResolveType(column.DBType, value)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%v", resolved), nil
		}
	case nil:
		return "null", nil
	default:
		err := fmt.Errorf("unknown value type %T for %s", value, column.DBType)
		return "", err
	}
}
