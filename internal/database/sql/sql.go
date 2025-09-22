package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/mysql"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/postgres"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/sqlite3"
	"github.com/ctrl-alt-boop/dribble/result"
)

const DefaultSelectLimit int = 10 // Just a safeguard

var SQLMethods = []Keyword{Select, Insert, Update, Delete, Pragma, Exec, Execute}

var _ database.SQL = &Executor{}

type Executor struct {
	db                   *sql.DB
	dialect              database.SQLDialect
	connectionProperties database.ConnectionProperties
}

// SetConnectionProperties implements database.SQL.
func (e *Executor) SetConnectionProperties(prop database.ConnectionProperties) {
	e.connectionProperties = prop
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
		return nil, fmt.Errorf("unknown or unsupported database dialect: %s", dialect)
	}
	if err != nil {
		return nil, err
	}
	return &Executor{
		dialect: dialect,
	}, nil
}

func (e *Executor) Open(_ context.Context) error {
	var connectionString strings.Builder
	e.dialect.ConnectionStringTemplate().Execute(&connectionString, e.connectionProperties)

	db, err := sql.Open(e.Dialect().Name(), connectionString.String())
	if err != nil {
		return err
	}
	e.db = db
	return nil
}

// Ping implements database.Executor.
func (e *Executor) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := e.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (e *Executor) Close(_ context.Context) error {
	return e.db.Close()
}

var ErrIntentTargetMismatch = errors.New("intent target does not match executor target")

// Execute implements database.Executor.
func (e *Executor) Execute(ctx context.Context, request database.Request) (any, error) {
	if err := e.Ping(ctx); err != nil {
		return nil, err
	}

	return e.execute(ctx, request)
}

// ExecuteWithHandler implements database.Executor.
func (e *Executor) ExecuteWithHandler(ctx context.Context, request database.Request, handler func(result any, err error)) {
	if err := e.Ping(ctx); err != nil {
		handler(nil, err)
	}

	handler(e.execute(ctx, request))
}

var ErrResultKindNotSupported = errors.New("result kind not supported")

func (e *Executor) execute(ctx context.Context, request database.Request) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var queryString string
	var queryArgs []any

	if request.IsPrefab() {
		prefabString, args, err := e.dialect.GetPrefab(request)
		if err != nil {
			return nil, err
		}
		queryString = prefabString
		queryArgs = args

	} else {
		requestString, args, err := e.dialect.RenderRequest(request)
		if err != nil {
			return nil, err
		}
		queryString = requestString
		queryArgs = args
	}
	kind := result.KindTable // FIXME: Temporary
	switch kind {
	case result.KindScalar:
		var scalar int
		row := e.db.QueryRowContext(ctx, queryString, queryArgs...)
		err := row.Scan(&scalar)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		return scalar, nil

	case result.KindRow:
		row := e.db.QueryRowContext(ctx, queryString, queryArgs...)
		return row, nil

	case result.KindList:
		rows, err := e.db.QueryContext(ctx, queryString, queryArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		return result.RowsToList(rows), nil

	case result.KindTable:
		rows, err := e.db.QueryContext(ctx, queryString, queryArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		if cols, _ := rows.Columns(); len(cols) == 1 {
			return result.RowsToList(rows), nil
		}
		return result.CreateDataTable(result.ParseRows(rows)), nil

	default:
		return nil, ErrResultKindNotSupported
	}
}

// Dialect implements database.SQL.
func (e *Executor) Dialect() database.SQLDialect {
	return e.dialect
}

// Request implements database.SQL.
func (e *Executor) Request(ctx context.Context, requests ...database.Request) (any, error) {
	return e.execute(ctx, requests[0])
}

// RequestWithHandler implements database.SQL.
func (e *Executor) RequestWithHandler(ctx context.Context, handler func(response database.Response, err error), requests ...database.Request) error {
	panic("unimplemented")
}

// Type implements database.SQL.
func (e *Executor) Type() database.Type {
	return database.TypeSQL
}
