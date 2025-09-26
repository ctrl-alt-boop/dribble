package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/datasource/sql/mysql"
	"github.com/ctrl-alt-boop/dribble/internal/datasource/sql/postgres"
	"github.com/ctrl-alt-boop/dribble/internal/datasource/sql/sqlite3"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/result"
)

var _ database.SQL = (*Executor)(nil)

type Executor struct {
	db              *sql.DB
	dialect         database.SQLDialect
	dataSourceNamer database.DataSourceNamer
}

func NewExecutor(dialectType database.SQLDialectType) (database.SQL, error) {
	var dialect database.SQLDialect
	var err error
	switch dialectType {
	case database.MySQL:
		dialect, err = mysql.NewMySQLDriver()
	case database.PostgreSQL:
		dialect, err = postgres.NewPostgresDriver()
	case database.SQLite3:
		dialect, err = sqlite3.NewSQLite3Driver()
	default:
		return nil, fmt.Errorf("unknown or unsupported database dialect: %s", dialectType)
	}
	if err != nil {
		return nil, err
	}
	return &Executor{
		dialect: dialect,
	}, nil
}

func New(dsn database.DataSourceNamer) (*Executor, error) {
	var dialect database.SQLDialect
	var err error
	switch dsn.Type() {
	case database.MySQL:
		dialect, err = mysql.NewMySQLDriver()
	case database.PostgreSQL:
		dialect, err = postgres.NewPostgresDriver()
	case database.SQLite3:
		dialect, err = sqlite3.NewSQLite3Driver()
	default:
		return nil, fmt.Errorf("unknown or unsupported database dialect: %s", dsn.Type())
	}
	if err != nil {
		return nil, err
	}
	return &Executor{
		dialect:         dialect,
		dataSourceNamer: dsn,
	}, nil
}

func (e *Executor) Open(ctx context.Context) error {
	connectionString := e.dataSourceNamer.DSN()

	fmt.Printf("connectionString: %+v\n", connectionString)

	if ctx.Err() != nil {
		return ctx.Err()
	}
	db, err := sql.Open(e.Dialect().Name(), connectionString)
	if err != nil {
		return err
	}
	e.db = db
	return nil
}

func (e *Executor) Ping(ctx context.Context) error {
	if e.db == nil {
		return errors.New("database connection is not open")
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := e.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (e *Executor) Close(_ context.Context) error {
	if e.db == nil {
		return nil // Already closed or never opened
	}
	return e.db.Close()
}

func (e *Executor) IsClosed() bool {
	return e.db == nil
}

// execute runs a single request against the database.
func (e *Executor) execute(ctx context.Context, req database.Request) (any, error) {
	if err := e.Ping(ctx); err != nil {
		return nil, err
	}

	intent, isIntent := req.(request.Intent)
	if !isIntent {
		if req.IsPrefab() {
			fmt.Printf("got prefab request: %T\n", req)
			queryString, queryArgs, err := e.dialect.GetPrefab(req)
			if err != nil {
				return nil, fmt.Errorf("failed to render prefab request: %w", err)
			}

			return e.executeRead(ctx, queryString, queryArgs)
		}

	}
	queryString, queryArgs, err := e.renderRequest(intent)
	if err != nil {
		return nil, fmt.Errorf("failed to render intent request: %w", err)
	}

	switch intent.Type {
	case database.Create, database.Update, database.Delete:
		return e.db.ExecContext(ctx, queryString, queryArgs...)
	case database.Read:
		return e.executeRead(ctx, queryString, queryArgs)
	default:
		// Fallback to Exec for unknown types, could also be an error.
		return e.db.ExecContext(ctx, queryString, queryArgs...)
	}
}

// renderRequest converts a database.Request into a query string and arguments.
func (e *Executor) renderRequest(intent request.Intent) (string, []any, error) {
	queryStringTemplate := e.dialect.GetTemplate(intent.Type)
	if queryStringTemplate == "" {
		return "", nil, fmt.Errorf("no template found for request type %v in dialect %s", intent.Type, e.dialect.Name())
	}

	tmpl, err := template.New("query").Parse(queryStringTemplate)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing query template: %w", err)
	}
	operation := intent.Operation
	var sb strings.Builder
	if err := tmpl.Execute(&sb, operation); err != nil {
		return "", nil, fmt.Errorf("error executing query template: %w", err)
	}
	return strings.TrimSpace(sb.String()), intent.Args, nil
}

func (e *Executor) executeRead(ctx context.Context, query string, args []any) (any, error) {
	fmt.Printf("queryString: %+v\n", query)
	rows, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// If the result has only one column, treat it as a list.
	if cols, _ := rows.Columns(); len(cols) == 1 {
		return result.RowsToList(rows), nil
	}
	columns, dataRows := result.ParseRows(rows)
	return result.NewTable(columns, dataRows), nil
}

func (e *Executor) Dialect() database.SQLDialect {
	return e.dialect
}

func (e *Executor) Request(ctx context.Context, request database.Request) (any, error) {
	// if len(requests) != 1 {
	// 	return nil, errors.New("SQL executor received multiple requests, but only supports one at a time")
	// }
	return e.execute(ctx, request)
}

// Type implements database.SQL.
func (e *Executor) Type() database.Type {
	return database.TypeSQL
}
