package datasource

import (
	"context"
	"text/template"
)

type (
	Capability   string
	Capabilities []Capability

	StorageType string

	Metadata struct {
		SourceType  DataSourceType
		StorageType StorageType
	}
)

// const (
// 	postgres = "postgres" // PostgreSQL wire protocol
// 	mysql    = "mysql"    // MySQL protocol
// 	mongodb  = "mongodb"  // MongoDB protocol
// 	redis    = "redis"    // Redis protocol
// 	sqlite3  = "sqlite3"  // SQLite (file-based, but protocol name)
// 	http     = "http"     // REST APIs
// 	grpc     = "grpc"     // gRPC services
// 	file     = "file"     // File system access
// )

const (
	SupportsJSON  Capability = "json"
	SupportsJSONB Capability = "jsonb"
	SupportsSQL   Capability = "sql"
	SupportsBSON  Capability = "json"
	SupportsCSV   Capability = "csv"

	IsFile     StorageType = "file"
	IsDatabase StorageType = "database"
	IsSQL      StorageType = "sql"
	IsNoSQL    StorageType = "nosql"
	IsGraph    StorageType = "graph"
	IsTime     StorageType = "timeseries"
)

func (c Capabilities) AsStrings() []string {
	list := make([]string, len(c))
	for i, cap := range c {
		list[i] = string(cap)
	}
	return list
}

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

type (
	Response interface {
		Code() int
		Message() string
	}

	Request interface {
		IsPrefab() bool
		ResponseOnSuccess() Response
		ResponseOnError() Response
	}

	Database interface { // Someday this should change to something like executor
		Type() Type

		Open(context.Context) error
		Ping(context.Context) error
		Close(context.Context) error
		IsClosed() bool

		Request(context.Context, Request) (any, error)
	}

	NoSQL interface {
		Database
		Client() NoSQLAdapter
	}

	SQLAdapter interface {
		Name() string

		ConnectionStringTemplate() *template.Template

		GetTemplate(RequestType) string
		GetPrefab(Request) (string, []any, error)

		RenderPlaceholder(index int) string
		IncreamentPlaceholder() string

		RenderTypeCast() string
		RenderCurrentTimestamp() string
		RenderValue(any) string
		QuoteRune() rune
		Quote(string) string

		// ResolveType is used when a database/sql driver doesn't resolve a type
		ResolveType(string, []byte) (any, error)
	}

	NoSQLAdapter interface { // FIXME: Translate a request to a NoSQL client method chain
		SetConnectionProperties(map[string]string) // FIXME: Temporary map for connection things

		Open(context.Context) error
		Ping(context.Context) error
		Close(context.Context) error

		Read(any)
		ReadMany(any)
		Create(any)
		Update(any)
		Delete(any)
	}
)
