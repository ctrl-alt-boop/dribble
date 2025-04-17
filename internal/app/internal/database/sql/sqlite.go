package sql

import (
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite3 struct{}

func CreateSQLite3Driver() (*SQLite3, error) {
	driver := &SQLite3{}
	err := driver.Load()
	if err != nil {
		return nil, err
	}
	return driver, nil
}

// ResolveType implements database.DbDriver.
func (d *SQLite3) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}

// IsFile implements database.DbDriver.
func (d *SQLite3) IsFile() bool {
	return false
}

// SupportsJsonResult implements database.DbDriver.
func (d *SQLite3) SupportsJsonResult() bool {
	return false
}

// ConnectionString implements database.GoolDbDriver.
func (d *SQLite3) ConnectionString(settings *connection.Settings) string {
	panic("unimplemented")
}

// Load implements database.GoolDbDriver.
func (d SQLite3) Load() error {
	panic("unimplemented")
}

// DatabasesQuery implements database.GoolDbDriver.
func (d *SQLite3) DatabasesQuery() string {
	panic("unimplemented")
}

func (d *SQLite3) DatabaseNameQuery() string {
	panic("unimplemented")
}

func (d *SQLite3) TableNamesQuery() string {
	return "SELECT name FROM sqlite_master WHERE type='table';"
}

func (d *SQLite3) CountQuery(table string) string {
	panic("unimplemented")
}

// SelectAllQuery implements database.GoolDbDriver.
func (d *SQLite3) SelectAllQuery(table string, opts query.Statement) string {
	panic("unimplemented")
}
