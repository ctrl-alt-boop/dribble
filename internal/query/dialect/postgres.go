package dialect

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/google/uuid"
)

type (
	Postgres struct {
	}
)

func (r *Postgres) Quote(value string) string {
	return "\"" + value + "\""
}

func (r *Postgres) QuoteRune() rune {
	return '"'
}

// Capabilities implements database.Dialect.
func (r *Postgres) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsJSON,
		database.SupportsJSONB,
	}
}

// GetTemplate implements database.Dialect.
func (r *Postgres) GetTemplate(queryType database.OperationType) string {
	switch queryType {
	case database.Read:
		return "" // DefaultSQLSelectTemplate
	case database.Create:
		return "" // DefaultSQLInsertTemplate
	case database.Update:
		return "" // DefaultSQLUpdateTemplate
	case database.Delete:
		return "" // DefaultSQLDeleteTemplate
	default:
		return ""
	}
}

// RenderCurrentTimestamp implements database.Dialect.
func (r *Postgres) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (r *Postgres) RenderPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// RenderTypeCast implements database.Dialect.
func (r *Postgres) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (r *Postgres) RenderValue(value any) string {
	panic("unimplemented")
}

func (r *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}
