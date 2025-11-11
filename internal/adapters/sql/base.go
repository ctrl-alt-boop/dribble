package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/internal/adapters"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/result"
)

type SQLExecutor interface {
	datasource.Executor
	GetTemplate(datasource.RequestType) string
	GetPrefab(datasource.Request) (string, []any, error)
}

type Base struct {
	adapters.BaseDatabase
	Self SQLExecutor
	DB   *sql.DB
	DSN  datasource.Namer
}

func (b *Base) Open(ctx context.Context) error {
	connectionString := b.DSN.DSN()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	db, err := sql.Open(b.Self.GoName(), connectionString)
	if err != nil {
		return err
	}
	b.DB = db
	return nil
}

func (b *Base) Ping(ctx context.Context) error {
	if b.DB == nil {
		return errors.New("database connection is not open")
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := b.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (b *Base) Close(_ context.Context) error {
	if b.DB == nil {
		return nil // Already closed or never opened
	}
	return b.DB.Close()
}

func (b *Base) IsClosed() bool {
	return b.DB == nil
}

// execute runs a single request against the database.
func (b *Base) execute(ctx context.Context, req datasource.Request) (any, error) {
	if err := b.Ping(ctx); err != nil {
		return nil, err
	}

	intent, isIntent := req.(request.Intent)
	if !isIntent {
		if req.IsPrefab() {
			fmt.Printf("got prefab request: %T\n", req)
			queryString, queryArgs, err := b.Self.GetPrefab(req)
			if err != nil {
				return nil, fmt.Errorf("failed to render prefab request: %w", err)
			}

			return b.executeRead(ctx, queryString, queryArgs)
		}
	}
	queryString, queryArgs, err := b.renderRequest(intent)
	if err != nil {
		return nil, fmt.Errorf("failed to render intent request: %w", err)
	}

	switch intent.Type {
	case datasource.Create, datasource.Update, datasource.Delete:
		return b.DB.ExecContext(ctx, queryString, queryArgs...)
	case datasource.Read:
		return b.executeRead(ctx, queryString, queryArgs)
	default:
		// Fallback to Exec for unknown types, could also be an error.
		return b.DB.ExecContext(ctx, queryString, queryArgs...)
	}
}

// renderRequest converts a database.Request into a query string and arguments.
func (b *Base) renderRequest(intent request.Intent) (string, []any, error) {
	queryStringTemplate := b.Self.GetTemplate(intent.Type)
	if queryStringTemplate == "" {
		return "", nil, fmt.Errorf("no template found for request type %v in dialect %s", intent.Type, b.Self.Name())
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

func (b *Base) executeRead(ctx context.Context, query string, args []any) (any, error) {
	fmt.Printf("queryString: %+v\n", query)
	rows, err := b.DB.QueryContext(ctx, query, args...)
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

func (b *Base) Request(ctx context.Context, request datasource.Request) (any, error) {
	// if len(requests) != 1 {
	// 	return nil, errors.New("SQL BaseSQL received multiple requests, but only supports one at a time")
	// }
	return b.execute(ctx, request)
}

// Type implements database.SQL.
func (b *Base) ExecutorType() datasource.ExecutorType {
	return datasource.ExecutorType("SQL")
}
