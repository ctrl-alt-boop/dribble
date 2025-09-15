package sql

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/sql"
)

var Prefabs *QueryPrefab = CreateQueryPrefabs()

type (
	QueryPrefab struct {
		Postgres  *Prefab
		MySQL     *Prefab
		SQLite3   *Prefab
		MsSQL     *Prefab
		Oracle    *Prefab
		SQLServer *Prefab
	}

	Prefab struct {
		CurrentDatabase *database.Intent
		Databases       *database.Intent
		Tables          *database.Intent
		Columns         func(tableName string) *database.Intent
	}
)

func CreateQueryPrefabs() *QueryPrefab {
	return &QueryPrefab{
		Postgres:  CreatePostgresPrefabs(),
		MySQL:     CreateMySQLPrefabs(),
		SQLite3:   CreateSQLitePrefabs(),
		MsSQL:     CreateMsSQLPrefabs(),
		Oracle:    CreateOraclePrefabs(),
		SQLServer: CreateSQLServerPrefabs(),
	}
}

func CreatePostgresPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT current_database()
		// Databases:       "SELECT datname FROM pg_database",
		// Tables:          "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'",
		// Columns:         "SELECT column_name FROM information_schema.columns WHERE table_name = $1",
		CurrentDatabase: sql.Select("current_database()").From("").ToIntent(),
		Databases:       sql.Select("datname").From("pg_database").ToIntent(),
		Tables:          sql.Select("table_name").From("information_schema.tables").Where(sql.Eq("table_schema", "public")).ToIntent(),
		Columns: func(tableName string) *database.Intent {
			return sql.Select("column_name").From("information_schema.columns").Where(sql.Eq("table_name", tableName)).ToIntent()
		},
	}
}

func CreateMySQLPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT DATABASE()",
		// Databases:       "SHOW DATABASES",
		// Tables:          "SHOW TABLES",
		// Columns:         "SHOW COLUMNS FROM $1",
		CurrentDatabase: sql.Select("DATABASE()").From("").ToIntent(),
		Databases:       sql.FromString("SHOW DATABASES"),
		Tables:          sql.FromString("SHOW TABLES"),
		Columns:         func(tableName string) *database.Intent { return sql.FromString("SHOW COLUMNS FROM $1", tableName) },
	}
}

func CreateSQLitePrefabs() *Prefab {
	return &Prefab{
		// Databases: "PRAGMA database_list",
		// Tables:    "SELECT name FROM sqlite_master WHERE type='table'",
		// Columns:   "PRAGMA table_info($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToIntent(),
		Databases:       sql.FromString("PRAGMA database_list"),
		Tables:          sql.Select("name").From("sqlite_master").Where(sql.Eq("type", "table")).ToIntent(),
		Columns: func(tableName string) *database.Intent {
			return sql.FromString("PRAGMA table_info($1)", tableName)
		},
	}
}

func CreateMsSQLPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT db_name()",
		// Databases:       "SELECT name FROM sys.databases",
		// Tables:          "SELECT name FROM sys.tables",
		// Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToIntent(),
		Databases:       sql.FromString("SELECT name FROM sys.databases"),
		Tables:          sql.FromString("SELECT name FROM sys.tables"),
		Columns: func(tableName string) *database.Intent {
			tableName = fmt.Sprintf("OBJECT_ID(%s)", tableName)
			return sql.Select("name").From("sys.columns").Where(sql.Eq("object_id", tableName)).ToIntent()
		},
	}
}

func CreateOraclePrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT sys_context('USERENV', 'CURRENT_SCHEMA') FROM dual",
		// Databases:       "SELECT owner FROM all_users",
		// Tables:          "SELECT table_name FROM all_tables",
		// Columns:         "SELECT column_name FROM all_tab_columns WHERE table_name = $1",
		CurrentDatabase: sql.Select("sys_context('USERENV', 'CURRENT_SCHEMA')").From("dual").ToIntent(),
		Databases:       sql.FromString("SELECT owner FROM all_users"),
		Tables:          sql.FromString("SELECT table_name FROM all_tables"),
		Columns: func(tableName string) *database.Intent {
			return sql.Select("column_name").From("all_tab_columns").Where(sql.Eq("table_name", tableName)).ToIntent()
		},
	}
}

func CreateSQLServerPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT db_name()",
		// Databases:       "SELECT name FROM sys.databases",
		// Tables:          "SELECT name FROM sys.tables",
		// Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToIntent(),
		Databases:       sql.FromString("SELECT name FROM sys.databases"),
		Tables:          sql.FromString("SELECT name FROM sys.tables"),
		Columns: func(tableName string) *database.Intent {
			tableName = fmt.Sprintf("OBJECT_ID(%s)", tableName)
			return sql.Select("name").From("sys.columns").Where(sql.Eq("object_id", tableName)).ToIntent()
		},
	}
}
