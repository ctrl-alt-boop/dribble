package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/result"
	_ "github.com/mattn/go-sqlite3"
)

var _ database.Driver = &SQLite3{}

type SQLite3 struct {
	Executor
	target *database.Target
}

func NewSQLite3Driver(target *database.Target) (*SQLite3, error) {
	if target.DriverName != "mysql" {
		return nil, fmt.Errorf("invalid driver name: %s", target.DriverName)
	}
	driver := &SQLite3{
		target: target,
	}
	return driver, nil
}

func (s *SQLite3) Close(_ context.Context) error {
	return s.DB.Close()
}

func (s *SQLite3) Open(_ context.Context) error {
	db, err := sql.Open("sqlite3", s.ConnectionString())
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *SQLite3) Dialect() database.Dialect {
	return s
}

func (s *SQLite3) Query(ctx context.Context, query *database.Intent) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	panic("unimplemented")
}

// func (s *SQLite3) execute(ctx context.Context, kind result.Kind, query string, queryArgs ...any) (any, error) {
// 	queryCtx, queryCancel := context.WithCancel(ctx)
// 	defer queryCancel()

// 	switch kind {
// 	case result.KindScalar:
// 		var scalar any
// 		row := s.DB.QueryRowContext(queryCtx, query, queryArgs...)
// 		err := row.Scan(&scalar)
// 		if err != nil {
// 			return queryError(fmt.Errorf("error executing query: %w", err))
// 		}
// 		return scalar, nil
// 	case result.KindList:
// 		rows, err := s.DB.QueryContext(queryCtx, query, queryArgs...)
// 		if err != nil {
// 			return queryError(fmt.Errorf("error executing query: %w", err))
// 		}
// 		defer rows.Close()

// 		return result.RowsToList(rows), nil
// 	case result.KindTable:
// 		rows, err := s.DB.QueryContext(queryCtx, query, queryArgs...)
// 		if err != nil {
// 			return queryError(fmt.Errorf("error executing query: %w", err))
// 		}
// 		defer rows.Close()

// 		return result.CreateDataTable(result.ParseRows(rows)), nil
// 	default:
// 		return queryError(fmt.Errorf("result kind not supported for postgres yet"))
// 	}
// }

func (s *SQLite3) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) (any, error) {
	switch prefabType {
	case database.PrefabTypeCurrentDatabase:
		return nil, fmt.Errorf("prefab type not supported for postgres")
	case database.PrefabTypeDatabases:
		return s.execute(ctx, result.KindList, Prefabs.Sqlite3.Databases, nil)
	case database.PrefabTypeTables:
		return s.execute(ctx, result.KindList, Prefabs.Sqlite3.Tables, nil)
	case database.PrefabTypeColumns:
		return s.execute(ctx, result.KindList, Prefabs.Sqlite3.Columns, args[0])
	default:
		return nil, fmt.Errorf("prefab type not supported for postgres")
	}
}

func (s *SQLite3) SetTarget(target *database.Target) {
	s.target = target
}

func (s *SQLite3) Target() *database.Target {
	return s.target
}

func (s *SQLite3) Ping(_ context.Context) error {
	return s.DB.Ping()
}

func (s *SQLite3) ConnectionString() string {
	return ""
}

func (d *SQLite3) Quote(value string) string {
	return "\"" + value + "\""
}

func (d *SQLite3) QuoteRune() rune {
	return '"'
}

// Capabilities implements database.Dialect.
func (s *SQLite3) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.IsFile,
	}
}

// GetTemplate implements database.Dialect.
func (s *SQLite3) GetTemplate(queryType database.OperationType) string {
	switch queryType {
	case database.Read:
		return DefaultSQLSelectTemplate
	case database.Create:
		return "" // DefaultSQLInsertTemplate
	case database.Update:
		return "" // DefaultSQLUpdateTemplate
	case database.Delete:
		return "" // DefaultSQLDeleteTemplate
	default:
		return ""
	}
}

// RenderCurrentTimestamp implements database.Dialect.
func (s *SQLite3) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (s *SQLite3) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (s *SQLite3) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (s *SQLite3) RenderValue(value any) string {
	panic("unimplemented")
}

func (s *SQLite3) ResolveType(dbType string, value []byte) (any, error) {
	// |go        | sqlite3           |
	// |----------|-------------------|
	// |nil       | null              |
	// |int       | integer           |
	// |int64     | integer           |
	// |float64   | float             |
	// |bool      | integer           |
	// |[]byte    | blob              |
	// |string    | text              |
	// |time.Time | timestamp/datetime|
	return string(value), nil
}
