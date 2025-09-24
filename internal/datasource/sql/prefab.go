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
		CurrentDatabase database.Request
		Databases       database.Request
		Tables          database.Request
		Columns         func(tableName string) database.Request
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
		CurrentDatabase: sql.Select("current_database()").From("").ToRequest(),
		Databases:       sql.Select("datname").From("pg_database").ToRequest(),
		Tables:          sql.Select("table_name").From("information_schema.tables").Where(sql.Eq("table_schema", "public")).ToRequest(),
		Columns: func(tableName string) database.Request {
			return sql.Select("column_name").From("information_schema.columns").Where(sql.Eq("table_name", tableName)).ToRequest()
		},
	}
}

func CreateMySQLPrefabs() *Prefab {
	databases, _ := sql.FromString("SHOW DATABASES")
	tables, _ := sql.FromString("SHOW TABLES")
	return &Prefab{
		// CurrentDatabase: "SELECT DATABASE()",
		// Databases:       "SHOW DATABASES",
		// Tables:          "SHOW TABLES",
		// Columns:         "SHOW COLUMNS FROM $1",
		CurrentDatabase: sql.Select("DATABASE()").From("").ToRequest(),
		Databases:       databases,
		Tables:          tables,
		Columns: func(tableName string) database.Request {
			columns, _ := sql.FromString("SHOW COLUMNS FROM $1", tableName)
			return columns
		},
	}
}

func CreateSQLitePrefabs() *Prefab {
	databases, _ := sql.FromString("PRAGMA database_list")
	return &Prefab{
		// Databases: "PRAGMA database_list",
		// Tables:    "SELECT name FROM sqlite_master WHERE type='table'",
		// Columns:   "PRAGMA table_info($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToRequest(),
		Databases:       databases,
		Tables:          sql.Select("name").From("sqlite_master").Where(sql.Eq("type", "table")).ToRequest(),
		Columns: func(tableName string) database.Request {
			columns, _ := sql.FromString("PRAGMA table_info($1)", tableName)
			return columns
		},
	}
}

func CreateMsSQLPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT db_name()",
		// Databases:       "SELECT name FROM sys.databases",
		// Tables:          "SELECT name FROM sys.tables",
		// Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToRequest(),
		Databases:       sql.Select("name").From("sys.databases").ToRequest(),
		Tables:          sql.Select("name").From("sys.tables").ToRequest(),
		Columns: func(tableName string) database.Request {
			tableName = fmt.Sprintf("OBJECT_ID(%s)", tableName)
			return sql.Select("name").From("sys.columns").Where(sql.Eq("object_id", tableName)).ToRequest()
		},
	}
}

func CreateOraclePrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT sys_context('USERENV', 'CURRENT_SCHEMA') FROM dual",
		// Databases:       "SELECT owner FROM all_users",
		// Tables:          "SELECT table_name FROM all_tables",
		// Columns:         "SELECT column_name FROM all_tab_columns WHERE table_name = $1",
		CurrentDatabase: sql.Select("sys_context('USERENV', 'CURRENT_SCHEMA')").From("dual").ToRequest(),
		Databases:       sql.Select("owner").From("all_users").ToRequest(),
		Tables:          sql.Select("table_name").From("all_tables").ToRequest(),
		Columns: func(tableName string) database.Request {
			return sql.Select("column_name").From("all_tab_columns").Where(sql.Eq("table_name", tableName)).ToRequest()
		},
	}
}

func CreateSQLServerPrefabs() *Prefab {
	return &Prefab{
		// CurrentDatabase: "SELECT db_name()",
		// Databases:       "SELECT name FROM sys.databases",
		// Tables:          "SELECT name FROM sys.tables",
		// Columns:         "SELECT name FROM sys.columns WHERE object_id = OBJECT_ID($1)",
		CurrentDatabase: sql.Select("db_name()").From("").ToRequest(),
		Databases:       sql.Select("name").From("sys.databases").ToRequest(),
		Tables:          sql.Select("name").From("sys.tables").ToRequest(),
		Columns: func(tableName string) database.Request {
			tableName = fmt.Sprintf("OBJECT_ID(%s)", tableName)
			return sql.Select("name").From("sys.columns").Where(sql.Eq("object_id", tableName)).ToRequest()
		},
	}
}
