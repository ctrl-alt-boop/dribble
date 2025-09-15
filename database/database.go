package database

import (
	"context"
	"reflect"
)

type OperationType int

const (
	Read OperationType = iota
	Create
	Update
	Delete
	Execute
	// Meta?
)

var OperationTypes = []OperationType{
	Read,
	Create,
	Update,
	Delete,
	Execute,
}

type PrefabType int

const (
	PrefabCurrentDatabase PrefabType = iota
	PrefabDatabases
	PrefabTables
	PrefabColumns
)

type Capabilities string

const (
	SupportsJson  Capabilities = "json"
	SupportsJsonB Capabilities = "jsonb"
	IsFile        Capabilities = "file"
	SupportsSQL   Capabilities = "sql"
)

type DatabaseType int

const (
	SQL DatabaseType = iota
	NoSQL
)

type (
	Dialect interface {
		GetTemplate(operationType OperationType) string
		GetPrefab(prefabType PrefabType) (string, bool)

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

	Driver interface {
		Dialect() Dialect

		ConnectionString(target *Target) string
		RenderIntent(intent *Intent) (string, error)
	}

	NoSQLClient interface {
		Dialect() Dialect
		Open(ctx context.Context, target *Target) error
		Close(ctx context.Context) error

		Read(any)
		ReadMany(any)
		Create(any)
		Update(any)
		Delete(any)
		// Execute()
	}

	Executor interface {
		SetTarget(target *Target)
		Target() *Target

		SetDriver(driver Driver)
		Driver() Driver

		Open(ctx context.Context) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Execute(ctx context.Context, intent *Intent) error
		ExecutePrefab(ctx context.Context, prefabType PrefabType, args ...any) error
		ExecuteAndHandle(ctx context.Context, intent *Intent, handler func(result any, err error)) error

		OnBefore(f func(intent *Intent, err error))
		OnAfter(f func(intent *Intent, err error))
		OnResult(f func(result any, err error))
	}

	Intent struct {
		Target     *Target
		TargetName string

		Type OperationType

		QueryType reflect.Type
		Operation any

		Args []any
	}

	IntentBatch []*Intent
)
