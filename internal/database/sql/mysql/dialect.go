package mysql

import "github.com/ctrl-alt-boop/dribble/database"

// GetPrefab implements database.Dialect.
func (m *MySQL) GetPrefab(prefabType database.PrefabType) (string, bool) {
	panic("unimplemented")
}

// GetTemplate implements database.Dialect.
func (m *MySQL) GetTemplate(operationType database.OperationType) string {
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
func (m *MySQL) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (m *MySQL) Quote(value string) string {
	return "`" + value + "`"
}

// QuoteRune implements database.Dialect.
func (m *MySQL) QuoteRune() rune {
	return '`'
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MySQL) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (m *MySQL) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (m *MySQL) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (m *MySQL) RenderValue(value any) string {
	panic("unimplemented")
}
