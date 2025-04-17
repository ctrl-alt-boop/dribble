package sql

import (
	"bytes"
	"fmt"
	"plugin"
	"text/template"

	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
	_ "github.com/go-sql-driver/mysql"
)

const (
	databaseNames = `SELECT SCHEMA_NAME
FROM information_schema.schemata
WHERE SCHEMA_NAME NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys');`

	simpleDatabaseNames = "SHOW DATABASES;"

	currentDatabase       = "SELECT DATABASE();"
	simpleCurrentDatabase = "SELECT SCHEMA();"

	tableNames = `SELECT table_name
FROM information_schema.tables
WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE';`
	simpleTableNames = "SHOW TABLES;"

	countTemplate = `SELECT COUNT(*) FROM %s;`

	selectTemplateMySQL = `SELECT {{.Columns}}
FROM {{.Table}}
WHERE {{.Where}}
ORDER BY {{.OrderBy}}
LIMIT {{.Limit}}
OFFSET {{.Offset}};`
)

type MySql struct{}

func CreateMySqlDriver() (*MySql, error) {
	driver := &MySql{}
	err := driver.Load()
	if err != nil {
		return nil, err
	}
	return driver, nil
}

// ResolveType implements database.DbDriver.
func (d *MySql) ResolveType(dbType string, value []byte) (any, error) {
	return string(value), nil
}

// IsFile implements database.DbDriver.
func (d *MySql) IsFile() bool {
	return true
}

// SupportsJsonResult implements database.DbDriver.
func (d *MySql) SupportsJsonResult() bool {
	return false
}

// Server=myServerAddress;Port=1234;Database=myDataBase;Uid=myUsername;Pwd=myPassword;
func (d *MySql) ConnectionString(settings *connection.Settings) string {
	connString := ""
	if settings.Port == 0 {
		settings.Port = 3306
	}
	if settings.DbName == "" {
		settings.DbName = "mysql"
	}

	connString += settings.Username
	connString += ":"
	connString += settings.Password
	connString += "@"
	connString += "tcp("
	connString += settings.Ip
	connString += ":"
	connString += fmt.Sprintf("%d", settings.Port)
	connString += ")/"
	connString += settings.DbName
	return connString
}

func (d MySql) Load() error {
	plug, err := plugin.Open("./plugins/mysql.so")
	if err != nil {
		return err
	}

	_, err = plug.Lookup("Loaded")
	if err != nil {
		return err
	}
	return nil
}

func (d *MySql) DatabasesQuery() string {
	return databaseNames
}

func (d *MySql) DatabaseNameQuery() string {
	return currentDatabase
}

func (d *MySql) TableNamesQuery() string {
	return tableNames
}

func (d *MySql) CountQuery(table string) string {
	return fmt.Sprintf(countTemplate, table)

}

// SelectAllQuery implements database.GoolDbDriver.
func (d *MySql) SelectAllQuery(table string, opts query.Statement) string {
	var buf bytes.Buffer
	template.Must(template.New("select").Parse(selectTemplateMySQL)).Execute(&buf, opts)
	query := buf.String()
	return query
}
