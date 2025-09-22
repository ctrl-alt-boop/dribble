package mysql

import (
	_ "embed"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

//go:embed templates/select.tmpl
var selectQueryTemplate string

// GetTemplate implements database.Dialect.
func (m *MySQL) GetTemplate(queryType database.RequestType) string {
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
func (m *MySQL) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (m *MySQL) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (m *MySQL) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MySQL) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (m *MySQL) RenderPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// RenderTypeCast implements database.Dialect.
func (m *MySQL) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (m *MySQL) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}
