package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/mysql"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/postgres"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/sqlite3"
	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/ctrl-alt-boop/dribble/target"
)

const DefaultSelectLimit int = 10 // Just a safeguard

var SQLMethods = []Keyword{Select, Insert, Update, Delete, Pragma, Exec, Execute}

func CreateClient(dialect database.SQLDialect) (database.SQL, error) {
	switch dialect {
	case database.MySQL:
		return mysql.NewMySQLDriver()
	case database.PostgreSQL:
		return postgres.NewPostgresDriver()
	case database.SQLite3:
		return sqlite3.NewSQLite3Driver()
	default:
		return nil, fmt.Errorf("unknown or unsupported database dialect: %s", dialect)
	}
}

type Executor struct {
	db     *sql.DB
	target *target.Target
	driver database.SQL
}

func NewExecutor(target *target.Target) *Executor {
	driver, err := CreateDriverFromTarget(target)
	if err != nil {
		panic(err)
	}
	return &Executor{
		target: target,
		driver: driver,
	}
}

func (e *Executor) Open(_ context.Context) error {
	connectionString := e.driver.ConnectionString(e.target)
	db, err := sql.Open(e.target.Dialect.String(), connectionString)
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

// Driver implements database.Executor.
func (e *Executor) Driver() database.SQL {
	return e.driver
}

var ErrIntentTargetMismatch = errors.New("intent target does not match executor target")

// Execute implements database.Executor.
func (e *Executor) Execute(ctx context.Context, request *database.Request) (any, error) {
	if intent.Target != e.target {
		return nil, ErrIntentTargetMismatch
	}

	if err := e.Ping(ctx); err != nil {
		return nil, err
	}

	return e.execute(ctx, intent)
}

// ExecutePrefab implements database.Executor.
func (e *Executor) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) (any, error) {
	if err := e.Ping(ctx); err != nil {
		return nil, err
	}

	prefab, ok := e.driver.Dialect().GetPrefab(prefabType)
	if !ok {
		return nil, fmt.Errorf("prefab type not found for driver: %s", e.target.Dialect)
	}

	if prefabType == database.PrefabTables {
		prefab = fmt.Sprintf(prefab, args...)
	}
	intent := database.NewReadIntent(e.target, prefab, args...)

	return e.execute(ctx, intent)
}

// ExecuteWithHandler implements database.Executor.
func (e *Executor) ExecuteWithHandler(ctx context.Context, intent *database.Intent, handler func(result any, err error)) {
	if intent.Target != e.target {
		handler(nil, ErrIntentTargetMismatch)
	}

	if err := e.Ping(ctx); err != nil {
		handler(nil, err)
	}

	handler(e.execute(ctx, intent))
}

// ExecuteWithChannel implements database.Executor.
func (e *Executor) ExecuteWithChannel(ctx context.Context, intent *database.Intent, eventChannel chan any) {
	if intent.Target != e.target {
		eventChannel <- ErrIntentTargetMismatch
	}

	if err := e.Ping(ctx); err != nil {
		eventChannel <- err
	}

	result, err := e.execute(ctx, intent)
	if err != nil {
		eventChannel <- err
	}
	eventChannel <- result
}

var ErrResultKindNotSupported = errors.New("result kind not supported")

func (e *Executor) execute(ctx context.Context, intent *database.Intent) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	queryString, err := e.driver.RenderIntent(intent)
	if err != nil {
		return nil, err
	}
	kind := result.DefaultOperationResults[intent.Type]
	switch kind {
	case result.KindScalar:
		var scalar int
		row := e.db.QueryRowContext(ctx, queryString, intent.Args...)
		err := row.Scan(&scalar)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		return scalar, nil

	case result.KindRow:
		row := e.db.QueryRowContext(ctx, queryString, intent.Args...)
		return row, nil

	case result.KindList:
		rows, err := e.db.QueryContext(ctx, queryString, intent.Args...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		return result.RowsToList(rows), nil

	case result.KindTable:
		rows, err := e.db.QueryContext(ctx, queryString, intent.Args...)
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
