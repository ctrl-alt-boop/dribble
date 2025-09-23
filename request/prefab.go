package request

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

var (
	_ database.Request = (*ReadDatabaseSchema)(nil)
	_ database.Request = (*ReadTableSchema)(nil)
	_ database.Request = (*ReadColumnSchema)(nil)
	_ database.Request = (*ReadDatabaseProperties)(nil)
	_ database.Request = (*ReadTableProperties)(nil)
	_ database.Request = (*ReadColumnProperties)(nil)
	_ database.Request = (*ReadDatabaseNames)(nil)
	_ database.Request = (*ReadTableNames)(nil)
	_ database.Request = (*ReadColumnNames)(nil)
	_ database.Request = (*ReadCount)(nil)
	_ database.Request = (*ReadAllCounts)(nil)
)

type prefab struct {
	successStatus Status
}

func (p prefab) IsPrefab() bool {
	return true
}

func (p prefab) Name() string {
	return fmt.Sprintf("%T", p)
}

// args = [database_name]
type ReadDatabaseSchema struct {
	prefab
	DatabaseName string
}

func NewReadDatabaseSchema(databaseName string) ReadDatabaseSchema {
	return ReadDatabaseSchema{
		prefab:       prefab{successStatus: SuccessReadDatabaseSchema},
		DatabaseName: databaseName,
	}
}

// args = [database name, table name]
type ReadTableSchema struct {
	prefab
	DatabaseName string
	TableName    string
}

func NewReadTableSchema(databaseName, tableName string) ReadTableSchema {
	return ReadTableSchema{
		prefab:       prefab{successStatus: SuccessReadTableSchema},
		DatabaseName: databaseName,
		TableName:    tableName,
	}
}

// args = [database name, table name, column name]
type ReadColumnSchema struct {
	prefab
	DatabaseName string
	TableName    string
	ColumnName   string
}

func NewReadColumnSchema(databaseName, tableName, columnName string) ReadColumnSchema {
	return ReadColumnSchema{
		prefab:       prefab{successStatus: SuccessReadColumnSchema},
		DatabaseName: databaseName,
		TableName:    tableName,
		ColumnName:   columnName,
	}
}

// args = [database name]
type ReadDatabaseProperties struct {
	prefab
	DatabaseName string
}

func NewReadDatabaseProperties(databaseName string) ReadDatabaseProperties {
	return ReadDatabaseProperties{
		prefab:       prefab{successStatus: SuccessReadDatabaseProperties},
		DatabaseName: databaseName,
	}
}

// args = [database name, table name]
type ReadTableProperties struct {
	prefab
	DatabaseName string
	TableName    string
}

func NewReadTableProperties(databaseName, tableName string) ReadTableProperties {
	return ReadTableProperties{
		prefab:       prefab{successStatus: SuccessReadTableProperties},
		DatabaseName: databaseName,
		TableName:    tableName,
	}
}

// args = [database name, table name, column name]
type ReadColumnProperties struct {
	prefab
	DatabaseName string
	TableName    string
	ColumnName   string
}

func NewReadColumnProperties(databaseName, tableName, columnName string) ReadColumnProperties {
	return ReadColumnProperties{
		prefab:       prefab{successStatus: SuccessReadColumnProperties},
		DatabaseName: databaseName,
		TableName:    tableName,
		ColumnName:   columnName,
	}
}

// args = [target name]
type ReadDatabaseNames struct {
	prefab
}

func NewReadDatabaseNames() ReadDatabaseNames {
	return ReadDatabaseNames{
		prefab: prefab{successStatus: SuccessReadDatabaseList},
	}
}

// args = [database name]
type ReadTableNames struct {
	prefab
	DatabaseName string
}

func NewReadTableNames() ReadTableNames {
	return ReadTableNames{
		prefab:       prefab{successStatus: SuccessReadDBTableList},
		DatabaseName: "",
	}
}

func NewReadDBTableNames(databaseName string) ReadTableNames {
	return ReadTableNames{
		prefab:       prefab{successStatus: SuccessReadDBTableList},
		DatabaseName: databaseName,
	}
}

// args = [database name, table name]
type ReadColumnNames struct {
	prefab
	DatabaseName string
	TableName    string
}

func NewReadColumnNames(databaseName, tableName string) ReadColumnNames {
	return ReadColumnNames{
		prefab:       prefab{successStatus: SuccessReadDBColumnList},
		DatabaseName: databaseName,
		TableName:    tableName,
	}
}

// args = [database name, table name]
type ReadCount struct {
	prefab
	DatabaseName string
	TableName    string
}

func NewReadCount(tableName string) ReadCount {
	return ReadCount{
		prefab:       prefab{successStatus: SuccessReadCount},
		DatabaseName: "",
		TableName:    tableName,
	}
}

func NewReadCountWithDB(databaseName, tableName string) ReadCount {
	return ReadCount{
		prefab:       prefab{successStatus: SuccessReadCount},
		DatabaseName: databaseName,
		TableName:    tableName,
	}
}

// args = [[database name, table name], [database name, table name], ...]
type ReadAllCounts struct {
	prefab
	DatabaseName string
	TableNames   []string
}

func NewReadAllCounts(databaseName string, tableNames []string) ReadAllCounts {
	return ReadAllCounts{
		prefab:       prefab{successStatus: SuccessReadCount},
		DatabaseName: databaseName,
		TableNames:   tableNames,
	}
}

func (p prefab) ResponseOnSuccess() database.Response {
	if p.successStatus == StatusUnknown {
		panic(fmt.Sprintf("prefab request of type %T created without using its constructor", p))
	}
	return Response{Status: p.successStatus}
}

func (p prefab) ResponseOnError() database.Response {
	if p.successStatus == StatusUnknown {
		panic(fmt.Sprintf("prefab request of type %T created without using its constructor", p))
	}
	return Response{Status: -p.successStatus}
}
