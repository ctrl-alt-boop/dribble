package sqlite3

import (
	_ "embed"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

//go:embed templates/select.tmpl
var selectQueryTemplate string

// GetTemplate implements database.Dialect.
func (s *SQLite3) GetTemplate(queryType database.RequestType) string {
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
func (s *SQLite3) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (s *SQLite3) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (s *SQLite3) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (s *SQLite3) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (s *SQLite3) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (s *SQLite3) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (s *SQLite3) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}
