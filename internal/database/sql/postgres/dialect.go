package postgres

import (
	_ "embed"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

//go:embed templates/select.tmpl
var selectQueryTemplate string

var _ database.Dialect = &Postgres{}

// GetTemplate implements database.Dialect.
func (p *Postgres) GetTemplate(queryType database.OperationType) string {
	switch queryType {
	case database.Read:
		return selectQueryTemplate
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

// IncreamentPlaceholder implements database.Dialect.
func (p *Postgres) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (p *Postgres) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (p *Postgres) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (p *Postgres) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (p *Postgres) RenderPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// RenderTypeCast implements database.Dialect.
func (p *Postgres) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (p *Postgres) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}
