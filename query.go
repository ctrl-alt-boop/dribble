package dribble

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/database/sql"
)

type (
	Query interface {
		ToSQL(dialect database.Dialect) (queryString string, params []any, err error)
		ToSQLFormatted(dialect database.Dialect) (queryString string, params []any, err error)
		Parameters() []any

		Method() sql.Method
	}

	Dependency struct {
		SourceId    int
		SouceDataId int

		TargetId     int
		TargetDataId int
	}

	Batch struct {
		Queries      []Query
		Dependencies []Dependency
	}
)

// func SQLStyleSelect() sql.SelectBuilder {
// 	return sql.SelectBuilder{}
// }

// func (c *Client) FetchDatabaseList() {
// 	if c.connections == nil {
// 		c.onEvent(DatabaseListFetchError, nil, ErrNoDatabaseConnection)
// 		return
// 	}
// 	list, err := c.connections.FetchDatabases()
// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(DatabaseListFetchError, nil, err)
// 		return
// 	}
// 	var targetList []*DBTarget
// 	for _, name := range list {
// 		dbSettings := c.target.Copy(WithDB(name))
// 		targetList = append(targetList, dbSettings)
// 	}

// 	c.onEvent(DatabaseListFetched, DatabaseListFetchData{Driver: c.connections.DriverName, Databases: targetList}, nil)
// }

// func (c *Client) FetchTableList(databaseName string) {
// 	if c.connections == nil {
// 		// logger.Warn(ErrNoDatabaseContext)
// 		c.onEvent(DBTableListFetchError, nil, ErrNoDatabaseConnection)
// 		return
// 	}
// 	c.Reconnect(connection.WithDB(databaseName))
// 	list, err := c.connections.FetchTableNames()
// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(DBTableListFetchError, nil, err)
// 		return
// 	}
// 	c.onEvent(DBTableListFetched, TableListFetchData{Database: databaseName, Tables: list}, nil)
// }

// func (c *Client) FetchTable(tableName string) {
// 	columns, rows, err := c.connections.FetchTable(tableName)
// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(TableFetchError, nil, err)
// 		return
// 	}
// 	tableData := result.CreateDataTable(columns, rows)
// 	c.onEvent(TableFetched, TableFetchData{TableName: tableName, Table: tableData}, nil)
// }

// func (c *Client) FetchTableColumns(tableName string) {
// 	columns, rows, err := c.connections.FetchTableColumns(tableName)
// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(TableFetchError, nil, err)
// 		return
// 	}
// 	tableData := result.CreateDataTable(columns, rows)
// 	c.onEvent(TableFetched, TableFetchData{TableName: tableName, Table: tableData}, nil)
// }

// func (c *Client) FetchCounts(tableNames []string) {
// 	counts, err := c.connections.FetchCounts(tableNames)
// 	c.onEvent(TableCountFetched, TableCountFetchData{Counts: counts, Errors: err}, nil)
// }

// func (c *Client) Query(queryStatment Query) {
// 	var err error

// 	queryString, params, err := queryStatment.ToSQL(c.executors.Driver)
// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(QueryExecuteError, nil, err)
// 		return
// 	}

// 	// method := query.Method()
// 	dimension := queryStatment.ShouldReturn()

// 	var res int
// 	var columns []result.Column
// 	var rows []result.Row

// 	switch dimension {
// 	case query.Result: // its fine
// 		res, err = c.executors.Execute(queryString, params...)
// 	case query.ResultScalar:
// 		fallthrough
// 	case query.ResultList:
// 		fallthrough
// 	case query.ResultTable:
// 		columns, rows, err = c.executors.Query(queryString, params...)
// 	}

// 	if err != nil {
// 		// logger.Warn(err)
// 		c.onEvent(QueryExecuteError, nil, err)
// 		return
// 	}
// 	queryData := QueryData{Query: queryStatment}
// 	switch dimension {
// 	case -1:
// 		queryData.Value = res
// 	case 0:
// 		queryData.Value = result.CreateDataScalar(rows)
// 	case 1:
// 		queryData.List = result.CreateDataList(rows)
// 	case 2:
// 		queryData.Table = result.CreateDataTable(columns, rows)
// 	}
// 	c.onEvent(QueryExecuted, queryData, nil)
// }
