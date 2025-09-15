package sqlite3

import "github.com/ctrl-alt-boop/dribble/database"

// GetPrefab implements database.Dialect.
func (s *SQLite3) GetPrefab(prefabType database.PrefabType) (string, bool) {
	panic("unimplemented")
}

// GetTemplate implements database.Dialect.
func (s *SQLite3) GetTemplate(operationType database.OperationType) string {
	switch operationType {
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

// IncreamentPlaceholder implements database.Dialect.
func (s *SQLite3) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (s *SQLite3) Quote(value string) string {
	return "`" + value + "`"
}

// QuoteRune implements database.Dialect.
func (s *SQLite3) QuoteRune() rune {
	return '`'
}

// RenderCurrentTimestamp implements database.Dialect.
func (s *SQLite3) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (s *SQLite3) RenderPlaceholder(index int) string {
	panic("unimplemented")
}

// RenderTypeCast implements database.Dialect.
func (s *SQLite3) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (s *SQLite3) RenderValue(value any) string {
	panic("unimplemented")
}
