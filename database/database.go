package database

import (
	"context"
	"text/template"
)

type Capabilities string

const (
	SupportsJSON    Capabilities = "json"
	SupportsJSONB   Capabilities = "jsonb"
	IsFile          Capabilities = "file"
	SupportsSQLLike Capabilities = "sql"
	SupportsBSON    Capabilities = "json"
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

var DBTypes = createBaseTree()

type (
	ResponseHandler func(Response)

	Response interface {
		Code() int
		Message() string
	}

	Request interface {
		IsPrefab() bool
		ResponseOnSuccess() Response
		ResponseOnError() Response
	}

	Database interface {
		Type() Type

		Open(context.Context) error
		Ping(context.Context) error
		Close(context.Context) error

		Request(context.Context, ...Request) (any, error)
		RequestWithHandler(context.Context, func(Response, error), ...Request) error
	}

	SQL interface {
		Database
		Dialect() SQLDialect
	}

	NoSQL interface {
		Database
		Client() NoSQLClient
	}

	SQLDialect interface {
		Name() string

		ConnectionStringTemplate() *template.Template

		GetTemplate(RequestType) string
		GetPrefab(Request) (string, []any, error)

		Capabilities() []Capabilities

		RenderPlaceholder(index int) string
		IncreamentPlaceholder() string

		RenderTypeCast() string
		RenderCurrentTimestamp() string
		RenderValue(any) string
		QuoteRune() rune
		Quote(string) string

		ResolveType(string, []byte) (any, error)
	}

	NoSQLClient interface {
		SetConnectionProperties(map[string]string) // FIXME: Temporary map for connection things

		Open(context.Context) error
		Ping(context.Context) error
		Close(context.Context) error

		// TODO: Translate a request to a NoSQL client method chain
		Read(any)
		ReadMany(any)
		Create(any)
		Update(any)
		Delete(any)
		// Execute()
	}
)
