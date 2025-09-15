package postgres

import (
	"github.com/ctrl-alt-boop/dribble/database"
)

const (
	PrefabCurrentDatabase = "SELECT current_database()"
	PrefabDatabases       = "SELECT datname FROM pg_database"
	PrefabTables          = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	PrefabColumns         = "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	PrefabColumnsFormat   = "SELECT column_name FROM information_schema.columns WHERE table_name = '%s'"
)

var Prefabs = map[database.PrefabType]string{
	database.PrefabCurrentDatabase: PrefabCurrentDatabase,
	database.PrefabDatabases:       PrefabDatabases,
	database.PrefabTables:          PrefabTables,
	database.PrefabColumns:         PrefabColumnsFormat,
}

// GetPrefabs implements database.Dialect.
func (p *Postgres) GetPrefab(prefabType database.PrefabType) (string, bool) {
	prefab, ok := Prefabs[prefabType]
	return prefab, ok
}
