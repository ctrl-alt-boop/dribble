package sql

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"plugin"
	"strings"
	"time"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

var _ database.Driver = &Postgres{}

type Postgres struct {
	DB     *sql.DB
	target *database.Target

	FetchLimit       int
	FetchLimitOffset int
}

func NewPostgresDriver(target *database.Target) (*Postgres, error) {
	if target.DriverName != "postgres" {
		return nil, fmt.Errorf("invalid driver name: %s", target.DriverName)
	}
	driver := &Postgres{
		target: target,
	}
	// err := driver.Load()
	// if err != nil {
	// 	return nil, err
	// }
	return driver, nil
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

func (p *Postgres) Close(_ context.Context) error {
	return p.DB.Close()
}

func (p *Postgres) Open(_ context.Context) error {
	db, err := sql.Open("postgres", p.ConnectionString())
	if err != nil {
		return err
	}
	p.DB = db
	return nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.DB.PingContext(context.TODO())
}

func (d *Postgres) ConnectionString() string {
	connString := ""
	connString += fmt.Sprintf("host=%s ", d.target.Ip)
	if d.target.Port == 0 {
		d.target.Port = 5432
	}
	connString += fmt.Sprintf("port=%d ", d.target.Port)
	connString += fmt.Sprintf("user=%s ", d.target.Username)
	connString += fmt.Sprintf("password=%s ", d.target.Password)
	if d.target.DBName == "" {
		d.target.DBName = "postgres"
	}
	connString += fmt.Sprintf("dbname=%s ", d.target.DBName)
	sslmode, ok := d.target.AdditionalSettings["sslmode"]
	if !ok {
		sslmode = "disable"
	}
	connString += fmt.Sprintf("sslmode=%s ", sslmode)

	return connString
}

func (d *Postgres) Dialect() database.Dialect {
	return d
}

func (d *Postgres) Query(query *database.QueryIntent) (any, error) {
	return d.QueryContext(context.Background(), query)
}

func (d *Postgres) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {
	pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
	defer pingCancel()

	err := d.DB.PingContext(pingCtx)
	if err != nil {
		return queryError(err)
	}

	queryString, err := d.intentToQueryString(query)
	if err != nil {
		return queryError(err)
	}

	queryCtx, queryCancel := context.WithCancel(ctx)
	defer queryCancel()

	switch d.resultKind(query) {
	case result.KindScalar:
		var scalar any
		row := d.DB.QueryRowContext(queryCtx, queryString, query.Args...)
		err := row.Scan(&scalar)
		if err != nil {
			return queryError(fmt.Errorf("error executing query: %w", err))
		}
		return scalar, nil
	case result.KindList:
		rows, err := d.DB.QueryContext(queryCtx, queryString, query.Args...)
		if err != nil {
			return queryError(fmt.Errorf("error executing query: %w", err))
		}
		defer rows.Close()

		return result.RowsToList(rows), nil
	case result.KindTable:
		rows, err := d.DB.QueryContext(queryCtx, queryString, query.Args...)
		if err != nil {
			return queryError(fmt.Errorf("error executing query: %w", err))
		}
		defer rows.Close()

		return result.CreateDataTable(result.ParseRows(rows)), nil
	default:
		return queryError(fmt.Errorf("result kind not supported for postgres yet"))
	}
}

func (d *Postgres) intentToQueryString(query *database.QueryIntent) (string, error) {
	var queryStringTemplate string
	switch query.QueryStyle {
	case database.SQL:
		queryStringTemplate = d.GetTemplate(query.Type)
	case database.NoSQL:
		return "", fmt.Errorf("NoSQL intent not yet supported for postgres")
	default:
		return "", fmt.Errorf("intent type not supported for postgres")
	}
	tmpl, err := template.New("query").Parse(queryStringTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing query template: %w", err)
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, query.SQLQuery)
	if err != nil {
		return "", fmt.Errorf("error executing query template: %w", err)
	}
	queryString := strings.TrimSpace(sb.String())
	return queryString, nil
}

func queryError(err error) (any, error) {
	return nil, err
}

func (d *Postgres) resultKind(query *database.QueryIntent) result.Kind {
	switch query.Type {
	case database.ReadQuery:
		if query.SQLQuery.IsCount {
			return result.KindScalar
		}
		if len(query.SQLQuery.Fields) == 1 && query.SQLQuery.Fields[0] != "*" {
			return result.KindList
		}
		return result.KindTable
	default:
		return result.KindNone
	}
}

func (d *Postgres) SetTarget(target *database.Target) {

	d.target = target
}

func (d *Postgres) Target() *database.Target {
	return d.target
}

func (d *Postgres) Quote(value string) string {
	return "\"" + value + "\""
}

func (d *Postgres) QuoteRune() rune {
	return '"'
}

// Capabilities implements database.Dialect.
func (d *Postgres) Capabilities() []database.DialectProperties {
	return []database.DialectProperties{
		database.SupportsJson,
		database.SupportsJsonB,
	}
}

// GetTemplate implements database.Dialect.
func (d *Postgres) GetTemplate(queryType database.QueryType) string {
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
func (d *Postgres) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (d *Postgres) RenderPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// RenderTypeCast implements database.Dialect.
func (d *Postgres) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (d *Postgres) RenderValue(value any) string {
	panic("unimplemented")
}

func (d *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}
