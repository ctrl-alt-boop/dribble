package firestore

import "github.com/ctrl-alt-boop/dribble/database"

// GetPrefab implements database.Dialect.
func (f *Firestore) GetPrefab(prefabType database.PrefabType) (string, bool) {
	panic("unimplemented")
}

// GetTemplate implements database.Dialect.
func (f *Firestore) GetTemplate(operationType database.OperationType) string {
	panic("unimplemented")
}

// IncreamentPlaceholder implements database.Dialect.
func (f *Firestore) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (f *Firestore) Quote(value string) string {
	panic("unimplemented")
}

// QuoteRune implements database.Dialect.
func (f *Firestore) QuoteRune() rune {
	panic("unimplemented")
}

// RenderCurrentTimestamp implements database.Dialect.
func (f *Firestore) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (f *Firestore) RenderPlaceholder(index int) string {
	panic("unimplemented")
}

// RenderTypeCast implements database.Dialect.
func (f *Firestore) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (f *Firestore) RenderValue(value any) string {
	panic("unimplemented")
}
