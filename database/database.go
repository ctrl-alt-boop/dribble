package database

import (
	"context"
)

type Capabilities string

const (
	SupportsJSON    Capabilities = "json"
	SupportsJSONB   Capabilities = "jsonb"
	IsFile          Capabilities = "file"
	SupportsSQLLike Capabilities = "sql"
	SupportsBSON    Capabilities = "json"
)

type Type int

const (
	TypeSQL Type = iota
	TypeNoSQL
	TypeGraph
	TypeTimeSeries

	TypeUndefined int = -1 // undefined
)

type RequestType int

const (
	NoOp RequestType = iota - 1
	Create
	Read
	Update
	Delete

	// Meta?
)

var RequestTypes = []RequestType{
	Create,
	Read,
	Update,
	Delete,
}

type SQLDialectType Type

const (
	PostgreSQL SQLDialectType = iota // postgres
	MySQL                            // mysql
	SQLite3                          // sqlite3

	NumSupportedSQLDialects
)

type NoSQLModelType Type

const (
	MongoDB   NoSQLModelType = iota // mongo
	Firestore                       // firestore
	Redis                           // redis

	NumSupportedNoSQLModels
)

var DBTypes = createBaseTree()

type (
	ResponseHandler func(response Response, err error)

	Response interface {
		Code() int
		Message() string
	}

	Request interface {
		ResponseOnSuccess() Response
		ResponseOnError() Response
	}

	Database interface {
		Type() Type

		Open(ctx context.Context) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Request(ctx context.Context, requests ...Request) (any, error)
		RequestWithHandler(ctx context.Context, handler func(response Response, err error), requests ...Request) error
	}

	SQL interface {
		Database
		Dialect() SQLDialect
		RenderRequest(request Request) (string, error)
		// ConnectionString(target *target.Target) string
		// RenderIntent(intent *Intent) (string, error)
	}

	NoSQL interface {
		Database
		Model() any // FIXME: Create model type

		// Open(ctx context.Context, target *target.Target) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Read(any)
		ReadMany(any)
		Create(any)
		Update(any)
		Delete(any)
		// Execute()
	}

	SQLDialect interface {
		SQL

		GetTemplate(operationType RequestType) string

		Capabilities() []Capabilities

		RenderPlaceholder(index int) string
		IncreamentPlaceholder() string

		RenderTypeCast() string
		RenderCurrentTimestamp() string
		RenderValue(value any) string
		QuoteRune() rune
		Quote(value string) string

		ResolveType(dbType string, value []byte) (any, error)
	}
)
