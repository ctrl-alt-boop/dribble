package mongodb

import "github.com/ctrl-alt-boop/dribble/database"

// GetTemplate implements database.Dialect.
func (m *MongoDB) GetTemplate(queryType database.OperationType) string {
	switch queryType {
	case database.Read:
		return "" // MongoDBSelectTemplate
	case database.Create:
		return ""
	case database.Update:
		return ""
	case database.Delete:
		return ""
	default:
		return ""
	}
}

// Quote implements database.Dialect.
func (m *MongoDB) GetPrefab(prefabType database.PrefabType) (string, bool) {
	panic("not implemented")
}

// Quote implements database.Dialect.
func (m *MongoDB) IncreamentPlaceholder() string {
	panic("not implemented")
}

// Quote implements database.Dialect.
func (m *MongoDB) Quote(value string) string {
	panic("unimplemented")
}

// QuoteRune implements database.Dialect.
func (m *MongoDB) QuoteRune() rune {
	panic("unimplemented")
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MongoDB) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (m *MongoDB) RenderPlaceholder(index int) string {
	panic("unimplemented")
}

// RenderTypeCast implements database.Dialect.
func (m *MongoDB) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (m *MongoDB) RenderValue(value any) string {
	panic("unimplemented")
}
