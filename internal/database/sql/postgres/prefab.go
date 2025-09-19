package postgres

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

const (
	PrefabCurrentDatabase = "SELECT current_database()"
	PrefabDatabases       = "SELECT datname FROM pg_database"
	PrefabTables          = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	PrefabColumns         = "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	PrefabColumnsFormat   = "SELECT column_name FROM information_schema.columns WHERE table_name = '%s'"
)

var Prefabs = map[database.Request]string{
	// PrefabCurrentDatabase: PrefabCurrentDatabase,
	request.ReadDatabaseNames: PrefabDatabases,
	PrefabTables:              PrefabTables,
	PrefabColumns:             PrefabColumnsFormat,
}

// GetPrefabs implements database.Dialect.
func (p *Postgres) GetPrefab(requestType database.RequestType) (string, bool) {
	prefab, ok := Prefabs[requestType]
	return prefab, ok
}
