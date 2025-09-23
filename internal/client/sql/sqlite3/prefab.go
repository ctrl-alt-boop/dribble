package sqlite3

const (
	PrefabDatabases = "PRAGMA database_list"
	PrefabTables    = "SELECT name FROM sqlite_master WHERE type='table'"
	PrefabColumns   = "PRAGMA table_info(?)"
)
