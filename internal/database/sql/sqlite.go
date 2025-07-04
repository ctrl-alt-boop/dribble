package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	_ "github.com/mattn/go-sqlite3"
)

var _ database.Driver = &SQLite3{}

type SQLite3 struct {
	DB     *sql.DB
	target *database.Target

	FetchLimit       int
	FetchLimitOffset int
}

func NewSQLite3Driver(target *database.Target) (*SQLite3, error) {
	if target.DriverName != "mysql" {
		return nil, fmt.Errorf("invalid driver name: %s", target.DriverName)
	}
	driver := &SQLite3{
		target: target,
	}
	err := driver.Load()
	if err != nil {
		return nil, err
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

func (s *SQLite3) Query(query *database.QueryIntent) (any, error) {
	return s.QueryContext(context.Background(), query)
}

func (s *SQLite3) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {

	panic("unimplemented")
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

func (s SQLite3) Load() error {
	return nil
}

func (d *SQLite3) Quote(value string) string {
	return "\"" + value + "\""
}

func (d *SQLite3) QuoteRune() rune {
	return '"'
}

// Capabilities implements database.Dialect.
func (s *SQLite3) Capabilities() []database.DialectProperties {
	return []database.DialectProperties{
		database.IsFile,
	}
}

// GetTemplate implements database.Dialect.
func (s *SQLite3) GetTemplate(queryType database.QueryType) string {
	switch queryType {
	case database.ReadQuery:
		return DefaultSQLSelectTemplate
	case database.CreateQuery:
		return "" // DefaultSQLInsertTemplate
	case database.UpdateQuery:
		return "" // DefaultSQLUpdateTemplate
	case database.DeleteQuery:
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
