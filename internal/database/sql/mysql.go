package sql

import (
	"context"
	"database/sql"
	"fmt"
	"plugin"

	"github.com/ctrl-alt-boop/dribble/database"
	_ "github.com/go-sql-driver/mysql"
)

var _ database.Driver = &MySql{}

type MySql struct {
	DB     *sql.DB
	target *database.Target

	FetchLimit       int
	FetchLimitOffset int
}

func NewMySqlDriver(target *database.Target) (*MySql, error) {
	if target.DriverName != "mysql" {
		return nil, fmt.Errorf("invalid driver name: %s", target.DriverName)
	}
	driver := &MySql{
		target: target,
	}
	// err := driver.Load()
	// if err != nil {
	// 	return nil, err
	// }
	return driver, nil
}

func (m MySql) Load() error {
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

func (m *MySql) Close(_ context.Context) error {
	return m.DB.Close()
}

func (m *MySql) Open(_ context.Context) error {
	db, err := sql.Open("mysql", m.ConnectionString())
	if err != nil {
		return err
	}
	m.DB = db
	return nil
}

func (m *MySql) Ping(_ context.Context) error {
	return m.DB.Ping()
}

// Server=myServerAddress;Port=1234;Database=myDataBase;Uid=myUsername;Pwd=myPassword;
func (m *MySql) ConnectionString() string {
	connString := ""
	if m.target.Port == 0 {
		m.target.Port = 3306
	}
	if m.target.DBName == "" {
		m.target.DBName = "mysql"
	}

	connString += m.target.Username
	connString += ":"
	connString += m.target.Password
	connString += "@"
	connString += "tcp("
	connString += m.target.Ip
	connString += ":"
	connString += fmt.Sprintf("%d", m.target.Port)
	connString += ")/"
	connString += m.target.DBName
	return connString
}

func (m *MySql) Dialect() database.Dialect {
	return m
}

func (m *MySql) Query(query *database.QueryIntent) (any, error) {
	return m.QueryContext(context.Background(), query)
}

func (m *MySql) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {
	panic("unimplemented")
}

func (m *MySql) SetTarget(target *database.Target) {
	panic("unimplemented")
}

func (m *MySql) Target() *database.Target {
	panic("unimplemented")
}

func (d *MySql) Quote(value string) string {
	return "`" + value + "`"
}

func (d *MySql) QuoteRune() rune {
	return '`'
}

// Capabilities implements database.Dialect.
func (m *MySql) Capabilities() []database.DialectProperties {
	return []database.DialectProperties{}
}

// GetTemplate implements database.Dialect.
func (m *MySql) GetTemplate(queryType database.QueryType) string {
	switch queryType {
	case database.ReadQuery:
		return DefaultSQLSelectTemplate
	case database.CreateQuery:
		return "" // DefaultSQLInsertTemplate
	case database.UpdateQuery:
		return "" // DefaultSQLUpdateTemplate
	case database.DeleteQuery:
		return "" // DefaultSQLDeleteTemplate
	default:
		return ""
	}
}

// RenderValue implements database.Dialect.
func (m *MySql) RenderValue(value any) string {
	panic("unimplemented")
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MySql) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (m *MySql) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (m *MySql) RenderTypeCast() string {
	panic("unimplemented")
}

func (m *MySql) ResolveType(dbType string, value []byte) (any, error) {
	return string(value), nil
}
