package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/result"
)

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
)

const DefaultSelectLimit int = 10 // Just a safeguard

var SQLMethods = []Method{MethodSelect, MethodInsert, MethodUpdate, MethodDelete}

var _ database.Executor = &Executor{}

type Executor struct {
	DB     *sql.DB
	target *database.Target
	driver database.Driver

	onResult func(result any, err error)
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
	e.DB = db
	return nil
}

func (e *Executor) Close(_ context.Context) error {
	return e.DB.Close()
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
	err := e.DB.PingContext(ctx)
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
	err := e.DB.PingContext(ctx)
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
	err := e.DB.PingContext(ctx)
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
		row := e.DB.QueryRowContext(ctx, query, queryArgs...)
		err := row.Scan(&scalar)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		return scalar, nil

	case result.KindRow:
		row := e.DB.QueryRowContext(ctx, query, queryArgs...)
		return row, nil

	case result.KindList:
		rows, err := e.DB.QueryContext(ctx, query, queryArgs...)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		return result.RowsToList(rows), nil

	case result.KindTable:
		rows, err := e.DB.QueryContext(ctx, query, queryArgs...)
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
