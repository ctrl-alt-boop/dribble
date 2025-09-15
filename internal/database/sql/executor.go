package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/mysql"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/postgres"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql/sqlite3"
	"github.com/ctrl-alt-boop/dribble/result"
)

const (
	PostgreSQL = "postgres"
	MySQL      = "mysql"
	SQLite     = "sqlite3"
)

var SupportedDrivers []string = []string{
	PostgreSQL,
	MySQL,
	SQLite,
}

var Defaults = map[string]*database.Target{
	PostgreSQL: {
		Type:       database.DBDriver,
		DriverName: PostgreSQL,
		Ip:         "127.0.0.1",
		Port:       5432,
		AdditionalSettings: map[string]string{
			"sslmode": "disable",
		},
	},
	MySQL: {
		Type:       database.DBDriver,
		DriverName: MySQL,
		Ip:         "127.0.0.1",
		Port:       3306,
	},
	SQLite: {
		Type:       database.DBDriver,
		DriverName: SQLite,
	},
}

const (
	Select            = "SELECT"
	Insert            = "INSERT"
	Update            = "UPDATE"
	Delete            = "DELETE"
	From              = "FROM"
	Where             = "WHERE"
	Set               = "SET"
	Values            = "VALUES"
	OrderBy           = "ORDER BY"
	Asc               = "ASC"
	Desc              = "DESC"
	Limit             = "LIMIT"
	Offset            = "OFFSET"
	GroupBy           = "GROUP BY"
	Having            = "HAVING"
	Join              = "JOIN"
	LeftJoin          = "LEFT JOIN"
	RightJoin         = "RIGHT JOIN"
	FullJoin          = "FULL JOIN"
	CrossJoin         = "CROSS JOIN"
	On                = "ON"
	As                = "AS"
	InnerJoin         = "INNER JOIN"
	OuterJoin         = "OUTER JOIN"
	Union             = "UNION"
	Intersect         = "INTERSECT"
	Except            = "EXCEPT"
	UnionAll          = "UNION ALL"
	IntersectAll      = "INTERSECT ALL"
	ExceptAll         = "EXCEPT ALL"
	Not               = "NOT"
	In                = "IN"
	Between           = "BETWEEN"
	And               = "AND"
	Or                = "OR"
	IsNull            = "IS NULL"
	IsNotNull         = "IS NOT NULL"
	IsTrue            = "IS TRUE"
	IsFalse           = "IS FALSE"
	IsUnknown         = "IS UNKNOWN"
	IsDistinctFrom    = "IS DISTINCT FROM"
	IsNotDistinctFrom = "IS NOT DISTINCT FROM"
	Like              = "LIKE"
	NotLike           = "NOT LIKE"
	Ilike             = "ILIKE"
	NotIlike          = "NOT ILIKE"
	Any               = "ANY"
	All               = "ALL"
	Exists            = "EXISTS"
	Some              = "SOME"
	Unique            = "UNIQUE"
	PrimaryKey        = "PRIMARY KEY"
	ForeignKey        = "FOREIGN KEY"
	Check             = "CHECK"
	Default           = "DEFAULT"
	Null              = "NULL"
	True              = "TRUE"
	False             = "FALSE"
	Unknown           = "UNKNOWN"
)

type Method string

func (s Method) String() string {
	return string(s)
}

const (
	MethodSelect  Method = "SELECT"
	MethodInsert  Method = "INSERT"
	MethodUpdate  Method = "UPDATE"
	MethodDelete  Method = "DELETE"
	MethodCall    Method = "CALL"
	MethodExec    Method = "EXEC"
	MethodExecute Method = "EXECUTE"
	MethodPragma  Method = "PRAGMA"
)

const DefaultSelectLimit int = 10 // Just a safeguard

var SQLMethods = []Method{MethodSelect, MethodInsert, MethodUpdate, MethodDelete}

func CreateDriverFromTarget(target *database.Target) (database.Driver, error) {
	switch target.DriverName {
	case MySQL:
		return mysql.NewMySQLDriver(target)
	case PostgreSQL:
		return postgres.NewPostgresDriver(target)
	case SQLite:
		return sqlite3.NewSQLite3Driver(target)
	default:
		return nil, fmt.Errorf("unknown or unsupported driver: %s", target.DriverName)
	}
}

var _ database.Executor = &Executor{}

type (
	IntentHandler func(intent *database.Intent, err error)
	ResultHandler func(result any, err error)
)

type Executor struct {
	db     *sql.DB
	target *database.Target
	driver database.Driver

	onBefore IntentHandler
	onAfter  IntentHandler
	onResult ResultHandler
}

func NewExecutor(target *database.Target) *Executor {
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
	driverName := e.target.DriverName
	if driverName == "" {
		return fmt.Errorf("no driver name provided")
	}
	connectionString := e.driver.ConnectionString(e.target)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	e.db = db
	return nil
}

func (e *Executor) Close(_ context.Context) error {
	return e.db.Close()
}

// Driver implements database.Executor.
func (e *Executor) Driver() database.Driver {
	return e.driver
}

// Execute implements database.Executor.
func (e *Executor) Execute(ctx context.Context, intent *database.Intent) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := e.db.PingContext(ctx)
	if err != nil {
		return err
	}

	intentString, err := e.driver.RenderIntent(intent)
	if err != nil {
		return err
	}

	go func() {
		e.onResult(e.execute(ctx, result.DefaultOperationResults[intent.Type], intentString, intent.Args...))
	}()

	return nil
}

// ExecuteAndHandle implements database.Executor.
func (e *Executor) ExecuteAndHandle(ctx context.Context, intent *database.Intent, handler func(result any, err error)) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := e.db.PingContext(ctx)
	if err != nil {
		return err
	}

	intentString, err := e.driver.RenderIntent(intent)
	if err != nil {
		return err
	}

	go func() {
		handler(e.execute(ctx, result.DefaultOperationResults[intent.Type], intentString, intent.Args...))
	}()

	return nil
}

// ExecutePrefab implements database.Executor.
func (e *Executor) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) error {
	prefab, ok := e.driver.Dialect().GetPrefab(prefabType)
	if !ok {
		return fmt.Errorf("prefab type not found")
	}

	go func() {
		e.onResult(e.execute(ctx, result.DefaultOperationResults[database.Read], prefab, args...))
	}()

	return nil
}

// Ping implements database.Executor.
func (e *Executor) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := e.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

// SetDriver implements database.Executor.
func (e *Executor) SetDriver(driver database.Driver) {
	e.driver = driver
}

// SetTarget implements database.Executor.
func (e *Executor) SetTarget(target *database.Target) {
	e.target = target
}

// Target implements database.Executor.
func (e *Executor) Target() *database.Target {
	return e.target
}

func (e *Executor) execute(ctx context.Context, kind result.Kind, query string, queryArgs ...any) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	switch kind {
	case result.KindScalar:
		var scalar int
		row := e.db.QueryRowContext(ctx, query, queryArgs...)
		err := row.Scan(&scalar)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		return scalar, nil

	case result.KindRow:
		row := e.db.QueryRowContext(ctx, query, queryArgs...)
		return row, nil

	case result.KindList:
		rows, err := e.db.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		return result.RowsToList(rows), nil

	case result.KindTable:
		rows, err := e.db.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		if cols, _ := rows.Columns(); len(cols) == 1 {
			return result.RowsToList(rows), nil
		}
		return result.CreateDataTable(result.ParseRows(rows)), nil

	default:
		return nil, fmt.Errorf("result kind not supported yet")
	}
}

// OnResult implements database.Executor.
func (e *Executor) OnBefore(f func(intent *database.Intent, err error)) {
	e.onBefore = f
}

// OnResult implements database.Executor.
func (e *Executor) OnAfter(f func(intent *database.Intent, err error)) {
	e.onAfter = f
}

// OnResult implements database.Executor.
func (e *Executor) OnResult(f func(result any, err error)) {
	e.onResult = f
}

const DefaultSQLSelectTemplate = `SELECT {{if .AsDistinct}}DISTINCT{{end}}
{{- range $i, $field := .Fields -}}
    {{if $i}}, {{end}}{{$field}}
{{- end}}
FROM {{.Table}}
{{- range .Joins}}
    {{.Type}} JOIN {{.Table}} ON {{.On}}
{{- end}}
{{- if .WhereClause}}
WHERE {{.WhereClause}}
{{- end}}
{{- if .GroupByClause}}
GROUP BY {{range $i, $field := .GroupByClause}}{{if $i}}, {{end}}{{$field}}{{end}}
{{- end}}
{{- if .HavingClause}}
HAVING {{.HavingClause}}
{{- end}}
{{- if .OrderByClause}}
ORDER BY {{range $i, $field := .OrderByClause}}{{if $i}}, {{end}}{{$field}}{{if .DescClause}} DESC{{end}}{{end}}
{{- end}}
{{- if .LimitClause}}
LIMIT {{.LimitClause}}
{{- end}}
{{- if .OffsetClause}}
OFFSET {{.OffsetClause}}
{{- end}}`
