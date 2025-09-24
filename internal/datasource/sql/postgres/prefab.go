package postgres

const (
	PrefabCurrentDatabase = "SELECT current_database()"
	PrefabDatabases       = "SELECT datname FROM pg_database"
	PrefabTables          = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	PrefabColumns         = "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	PrefabColumnsFormat   = "SELECT column_name FROM information_schema.columns WHERE table_name = '%s'"
	PrefabCountFormat     = "SELECT COUNT(*) FROM %s"
	PrefabCountDBFormat   = "SELECT COUNT(*) FROM %s.%s"
)
