package sql

import (
	"fmt"
	"plugin"
	"strings"

	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	"github.com/google/uuid"
)

const (
	selectTemplate = `SELECT {{.Columns}} FROM {{.Table}} WHERE {{.Where}} ORDER BY {{.OrderBy}} LIMIT {{.Limit}} OFFSET {{.Offset}}`
	insertTemplate = `INSERT INTO {{.Table}} ({{.Columns}}) VALUES ({{.Values}})`
	updateTemplate = `UPDATE {{.Table}} SET {{.Set}} WHERE {{.Where}}`
	deleteTemplate = `DELETE FROM {{.Table}} WHERE {{.Where}}`
)

type Postgres struct{}

func CreatePostgresDriver() (*Postgres, error) {
	driver := &Postgres{}
	err := driver.Load()
	if err != nil {
		return nil, err
	}
	return driver, nil
}

// ResolveType implements database.DbDriver.
func (d *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}

// IsFile implements database.DbDriver.
func (d *Postgres) IsFile() bool {
	return false
}

// SupportsJsonResult implements database.DbDriver.
func (d *Postgres) SupportsJsonResult() bool {
	return true
}

func (d *Postgres) ConnectionString(settings *connection.Settings) string {
	connString := ""
	connString += fmt.Sprintf("host=%s ", settings.Ip)
	if settings.Port == 0 {
		settings.Port = 5432
	}
	connString += fmt.Sprintf("port=%d ", settings.Port)
	connString += fmt.Sprintf("user=%s ", settings.Username)
	connString += fmt.Sprintf("password=%s ", settings.Password)
	if settings.DbName == "" {
		settings.DbName = "postgres"
	}
	connString += fmt.Sprintf("dbname=%s ", settings.DbName)
	sslmode, ok := settings.AdditionalSettings["sslmode"]
	if !ok {
		sslmode = "disable"
	}
	connString += fmt.Sprintf("sslmode=%s ", sslmode)

	return connString
}

func (d Postgres) Load() error {
	plug, err := plugin.Open("./plugins/postgres.so")
	if err != nil {
		return err
	}

	_, err = plug.Lookup("Loaded")
	if err != nil {
		return err
	}
	return nil
}

func (d *Postgres) DatabasesQuery() string {
	return "SELECT datname FROM pg_database WHERE datistemplate = false"
}

func (d *Postgres) DatabaseNameQuery() string {
	return "SELECT current_database()"
}

func (d *Postgres) TableNamesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE'"
}

func (d *Postgres) CountQuery(table string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
}

func (d *Postgres) SelectAllQuery(table string, opts query.Statement) string {
	queryOptsString := buildQueryOptions(opts)
	return fmt.Sprintf("SELECT * FROM %s%s", table, queryOptsString)
}

func (d *Postgres) TablePropertiesQuery(table string) string {
	return ""
}

func (d *Postgres) TableIndexesQuery(table string) string {
	return ""
}

func (d *Postgres) TableSizeQuery(table string) string {
	return ""
}

func (d *Postgres) DatabasePropertiesQuery(database string) string {
	return ""
}

func (d *Postgres) DatabaseSizeQuery(database string) string {
	return ""
}

// internal
const (
	jsonRow query.SqlMethod = "SELECT row_to_json(t) FROM %s t"
	jsonAgg query.SqlMethod = "SELECT json_agg(row_to_json(t)) FROM %s t"
)

func (d *Postgres) WithJsonRow(statement *query.Statement) {
	statement.Method = jsonRow
}

func (d *Postgres) WithJsonAgg(statement *query.Statement) {
	statement.Method = jsonAgg
}

func buildQueryOptions(opts query.Statement) string {
	sb := strings.Builder{}

	// TODO: ADD WHERE THING HERE

	if opts.OrderBy.Column != "" {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(opts.OrderBy.Column)
		if opts.OrderBy.Desc {
			sb.WriteString(" DESC")
		} else {
			sb.WriteString(" ASC")
		}
	}

	if opts.Limit > 0 {
		sb.WriteString(fmt.Sprint(" LIMIT ", opts.Limit))
	}
	if opts.Offset > 0 {
		sb.WriteString(fmt.Sprint(" OFFSET ", opts.Offset))
	}

	return sb.String()
}
