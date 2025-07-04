package dribble

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

var DefaultFetchLimit = 10

type QueryExecuter struct {
	database.Driver
	DriverName DriverName

	FetchLimit int

	onQueryExecuted func(query string, err error)
}

func createQueryExecuter(target *database.Target) (*QueryExecuter, error) {
	driver, err := CreateDriverFromTarget(target)
	if err != nil {
		return nil, fmt.Errorf("error creating handler for driver: %w", err)
	}

	connection := &QueryExecuter{
		Driver:     driver,
		DriverName: target.DriverName,

		FetchLimit: DefaultFetchLimit,

		onQueryExecuted: func(query string, err error) {},
	}
	err = connection.Open(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	return connection, nil
}

func (e *QueryExecuter) OnQueryExecuted(f func(query string, err error)) {
	e.onQueryExecuted = f
}

func (e *QueryExecuter) VerifyConnection() error {
	err := e.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("error trying to ping database: %w", err)
	}
	return nil
}

// func (e *QueryExecuter) Query(query *database.QueryIntent) (any, error) {
// 	return e.QueryContext(context.Background(), query)
// }

// func (e *QueryExecuter) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {

// 	return e.dialect.GetTemplate(query.Type)
// }

// func (c *Connection) FetchDatabases() ([]string, error) {
// 	var databases []string
// 	rows, err := c.DB.Query(c.Driver.DatabasesQuery())
// 	c.onQueryExecuted(c.Driver.DatabasesQuery(), err)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching database name list: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var databaseName string
// 		if err := rows.Scan(&databaseName); err != nil {
// 			// logger.Warn(err)
// 		}
// 		databases = append(databases, databaseName)
// 	}

// 	return databases, nil
// }

// func (c *Connection) FetchDatabaseName() (string, error) {
// 	var dbName string
// 	queryString := c.Driver.DatabaseNameQuery()
// 	err := c.DB.QueryRow(queryString).Scan(&dbName)
// 	c.onQueryExecuted(queryString, err)
// 	if err != nil {
// 		return "", fmt.Errorf("error fetching database name: %w", err)
// 	}
// 	return dbName, nil
// }

// func (c *Connection) FetchTableNames() ([]string, error) {
// 	var tableNames []string
// 	rows, err := c.DB.Query(c.Driver.TableNamesQuery())
// 	c.onQueryExecuted(c.Driver.TableNamesQuery(), err)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching table names: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var tableName string
// 		if err := rows.Scan(&tableName); err != nil {
// 			// logger.Warn(err)
// 		}
// 		tableNames = append(tableNames, tableName)
// 	}

// 	// logger.Infof("Fetched %d tables %v", len(tableNames), tableNames)
// 	return tableNames, nil
// }

// func (c *Connection) FetchCount(tableName string) (int, error) {
// 	var count int
// 	err := c.DB.QueryRow(c.Driver.CountQuery(tableName)).Scan(&count)
// 	c.onQueryExecuted(c.Driver.CountQuery(tableName), err)
// 	if err != nil {
// 		return 0, fmt.Errorf("error fetching count: %w", err)
// 	}
// 	return count, nil
// }

// func (c *Connection) FetchCounts(tableNames []string) (map[string]int, map[string]error) {
// 	counts := make(map[string]int)
// 	errors := make(map[string]error)
// 	for _, table := range tableNames {
// 		count, err := c.FetchCount(table)
// 		if err != nil {
// 			// logger.Error(err)
// 			errors[table] = err
// 			continue
// 		}
// 		counts[table] = count
// 	}
// 	return counts, errors
// }

// func (c *Connection) Execute(query string, args ...any) (int, error) {
// 	res, err := c.DB.Exec(query, args...)
// 	c.onQueryExecuted(query, err)
// 	if err != nil {
// 		// logger.Warn(err)
// 		return 0, err
// 	}
// 	aff, err := res.RowsAffected()
// 	return int(aff), err
// }

// func (c *Connection) Query(query string, args ...any) ([]result.Column, []result.Row, error) {
// 	dbRows, err := c.DB.Query(query, args...)

// 	c.onQueryExecuted(query, err)
// 	if err != nil {
// 		// logger.Warn(err)
// 		return nil, nil, err
// 	}
// 	defer dbRows.Close()

// 	return result.ParseRows(c.Driver, dbRows)
// }

// func (c *Connection) FetchTable(tableName string) ([]result.Column, []result.Row, error) { // context.FetchLimitOffset += context.FetchLimit
// 	// selectQuery := query.SelectAll().From(tableName).Limit(context.FetchLimit).Offset(context.FetchLimitOffset)
// 	// dbRows, err := context.Query("")
// 	// if err != nil {
// 	// logger.Warn(err)
// 	// 	return nil, nil, err
// 	// }
// 	// defer dbRows.Close()

// 	// return ParseRows(context.Driver, dbRows)
// 	return nil, nil, nil
// }

// func (c *Connection) FetchTableColumns(tableName string) ([]result.Column, []result.Row, error) {
// 	dbRows, err := c.DB.Query(c.Driver.TableColumnsPropertiesQuery(tableName))
// 	c.onQueryExecuted(c.Driver.TableColumnsPropertiesQuery(tableName), err)
// 	if err != nil {
// 		// logger.Warn(err)
// 		return nil, nil, err
// 	}
// 	defer dbRows.Close()

// 	return result.ParseRows(c.Driver, dbRows)
// }
