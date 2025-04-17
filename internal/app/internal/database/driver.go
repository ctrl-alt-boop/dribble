package database

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database/sql"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

type DriverName = string

const (
	DriverPostgreSQL DriverName = "postgres"
	DriverMySql      DriverName = "mysql"
	DriverSQLite     DriverName = "sqlite3"
)

type DbDriver interface {
	data.Resolver
	SupportsJsonResult() bool
	IsFile() bool

	Load() error
	ConnectionString(settings *connection.Settings) string

	// CustomQuery(opts query.Statement)

	DatabasesQuery() string
	DatabaseNameQuery() string
	TableNamesQuery() string
	CountQuery(table string) string
	SelectAllQuery(table string, opts query.Statement) string
	// ResolveDatabaseType(dbType string, value []byte) (any, error)
}

func createDriver(name DriverName) (DbDriver, error) {
	switch name {
	case DriverMySql:
		return sql.CreateMySqlDriver()
	case DriverPostgreSQL:
		return sql.CreatePostgresDriver()
	case DriverSQLite:
		return sql.CreateSQLite3Driver()
	default:
		return nil, fmt.Errorf("unknown driver name: %s", name)
	}
}
