package database

import "context"

type QueryType int

const (
	ReadQuery QueryType = iota
	CreateQuery
	UpdateQuery
	DeleteQuery
	ExecuteQuery
	// Meta?
)

type DialectProperties string

const (
	SupportsJson  DialectProperties = "json"
	SupportsJsonB DialectProperties = "jsonb"
	IsFile        DialectProperties = "file"
)

type QueryStyle int

const (
	SQL QueryStyle = iota
	NoSQL
)

type (
	Dialect interface {
		GetTemplate(queryType QueryType) string
		Capabilities() []DialectProperties

		RenderPlaceholder(index int) string
		RenderTypeCast() string
		RenderCurrentTimestamp() string
		RenderValue(value any) string
		Quote(value string) string
		QuoteRune() rune

		ResolveType(dbType string, value []byte) (any, error)
	}

	Driver interface {
		SetTarget(target *Target)
		Target() *Target
		Dialect() Dialect

		Open(ctx context.Context) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Query(query *QueryIntent) (any, error)
		QueryContext(ctx context.Context, query *QueryIntent) (any, error)
	}

	QueryIntent struct {
		Type       QueryType
		QueryStyle QueryStyle
		TargetName string
		SQLQuery   *SQLSelectQuery
		NoSQLQuery *NoSQLSelectQuery

		Args []any
	}
)

type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

type (
	NoSQLSelectQuery struct {
		Collection string

		ConditionsClause string
		Conditions       Exprs
		args             []any

		LimitClause  *int
		OffsetClause *int
	}
	SQLSelectQuery struct {
		AsDistinct bool
		IsCount    bool

		Fields []string
		Table  string

		Joins []joinClause

		WhereClause string
		args        []any

		GroupByClause []string
		HavingClause  string
		OrderByClause []orderByClause

		LimitClause  *int
		OffsetClause *int
	}
)
