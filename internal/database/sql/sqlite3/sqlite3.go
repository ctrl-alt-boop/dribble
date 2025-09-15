package sqlite3

import (
	"github.com/ctrl-alt-boop/dribble/database"
	_ "github.com/mattn/go-sqlite3"
)

var _ database.Driver = &SQLite3{}
var _ database.Dialect = &SQLite3{}

type SQLite3 struct{}

func NewSQLite3Driver(target *database.Target) (*SQLite3, error) {
	driver := &SQLite3{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (s *SQLite3) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.IsFile,
	}
}

// ConnectionString implements database.Driver.
func (s *SQLite3) ConnectionString(target *database.Target) string {
	panic("unimplemented")
}

// Dialect implements database.Driver.
func (s *SQLite3) Dialect() database.Dialect {
	return s
}

// RenderIntent implements database.Driver.
func (s *SQLite3) RenderIntent(intent *database.Intent) (string, error) {
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
