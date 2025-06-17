package gooldb

import (
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
)

func (gool *GoolDb) FetchDatabaseList() {
	if gool.databaseContext == nil {
		gool.onEvent(DatabaseListFetchError, nil, ErrNoDatabaseContext)
		return
	}
	list, err := gool.databaseContext.FetchDatabases()
	if err != nil {
		logger.Warn(err)
		gool.onEvent(DatabaseListFetchError, nil, err)
		return
	}
	var settingsList []*connection.Settings
	for _, name := range list {
		dbSettings := gool.settings.Copy(connection.WithDb(name))
		settingsList = append(settingsList, dbSettings)
	}

	gool.onEvent(DatabaseListFetched, DatabaseListFetchData{Driver: gool.databaseContext.DriverName, Databases: settingsList}, nil)
}

func (gool *GoolDb) FetchTableList(databaseName string) {
	if gool.databaseContext == nil {
		logger.Warn(ErrNoDatabaseContext)
		gool.onEvent(DBTableListFetchError, nil, ErrNoDatabaseContext)
		return
	}
	gool.Reconnect(connection.WithDb(databaseName))
	list, err := gool.databaseContext.FetchTableNames()
	if err != nil {
		logger.Warn(err)
		gool.onEvent(DBTableListFetchError, nil, err)
		return
	}
	gool.onEvent(DBTableListFetched, DBTableListFetchData{Database: databaseName, Tables: list}, nil)
}

func (gool *GoolDb) FetchTable(tableName string) {
	columns, rows, err := gool.databaseContext.FetchTable(tableName)
	if err != nil {
		logger.Warn(err)
		gool.onEvent(TableFetchError, nil, err)
		return
	}
	tableData := data.CreateDataTable(columns, rows)
	gool.onEvent(TableFetched, TableFetchData{TableName: tableName, Table: tableData}, nil)
}

func (gool *GoolDb) FetchCounts(tableNames []string) {
	counts, err := gool.databaseContext.FetchCounts(tableNames)
	gool.onEvent(TableCountFetched, TableCountFetchData{Counts: counts, Errors: err}, nil)
}
