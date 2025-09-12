package sql

var Prefabs *QueryPrefab = CreateQueryPrefabs()

//"mongodb"

type (
	QueryPrefab struct {
		Postgres  *Prefab
		MySql     *Prefab
		Sqlite3   *Prefab
		Mssql     *Prefab
		Oracle    *Prefab
		SqlServer *Prefab
	}

	Prefab struct {
		CurrentDatabase string
		Databases       string
		Tables          string
		Columns         string
	}
)

func CreateQueryPrefabs() *QueryPrefab {
	return &QueryPrefab{
		Postgres:  CreatePostgresPrefabs(),
		MySql:     CreateMySqlPrefabs(),
		Sqlite3:   CreateSqlitePrefabs(),
		Mssql:     CreateMssqlPrefabs(),
		Oracle:    CreateOraclePrefabs(),
		SqlServer: CreateSqlServerPrefabs(),
	}
}

func CreatePostgresPrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "SELECT current_database()",
		Databases:       "SELECT datname FROM pg_database",
		Tables:          "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'",
		Columns:         "SELECT column_name FROM information_schema.columns WHERE table_name = $1",
	}
}

func CreateMySqlPrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "SELECT DATABASE()",
		Databases:       "SHOW DATABASES",
		Tables:          "SHOW TABLES",
		Columns:         "SHOW COLUMNS FROM $1",
	}
}

func CreateSqlitePrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "PRAGMA database_list",
		Databases:       "SELECT name FROM sqlite_master WHERE type='table'",
		Tables:          "SELECT name FROM sqlite_master WHERE type='table'",
		Columns:         "PRAGMA table_info($1)",
	}
}

func CreateMssqlPrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "SELECT db_name()",
		Databases:       "SELECT name FROM sys.databases",
		Tables:          "SELECT name FROM sys.tables",
		Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
	}
}

func CreateOraclePrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "SELECT sys_context('USERENV', 'CURRENT_SCHEMA') FROM dual",
		Databases:       "SELECT owner FROM all_users",
		Tables:          "SELECT table_name FROM all_tables",
		Columns:         "SELECT column_name FROM all_tab_columns WHERE table_name = $1",
	}
}

func CreateSqlServerPrefabs() *Prefab {
	return &Prefab{
		CurrentDatabase: "SELECT db_name()",
		Databases:       "SELECT name FROM sys.databases",
		Tables:          "SELECT name FROM sys.tables",
		Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
	}
}
