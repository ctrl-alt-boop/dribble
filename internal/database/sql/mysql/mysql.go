package mysql

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.DBTypes.Register("SQL", "mysql")
}

var _ database.SQL = &MySQL{}
var _ database.SQLDialect = &MySQL{}

type MySQL struct{}

func NewMySQLDriver() (*MySQL, error) {
	driver := &MySQL{}
	return driver, nil
}

// // Capabilities implements database.Dialect.
func (m *MySQL) Capabilities() []database.Capabilities {
	return []database.Capabilities{}
}

// ConnectionString implements database.Driver.
// Server=myServerAddress;Port=1234;Database=myDataBase;Uid=myUsername;Pwd=myPassword;
func (m *MySQL) ConnectionString(target *target.Target) string {
	connString := ""
	if target.Port == 0 {
		target.Port = 3306
	}
	if target.DBName == "" {
		target.DBName = "mysql"
	}

	connString += target.Username
	connString += ":"
	connString += target.Password
	connString += "@"
	connString += "tcp("
	connString += target.Ip
	connString += ":"
	connString += fmt.Sprintf("%d", target.Port)
	connString += ")/"
	connString += target.DBName
	return connString
}

// Dialect implements database.Driver.
func (m *MySQL) Dialect() database.SQLDialect {
	return m
}

// RenderRequest implements database.Driver.
func (m *MySQL) RenderRequest(intent *database.Request) (string, error) {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (m *MySQL) ResolveType(dbType string, value []byte) (any, error) {
	return string(value), nil
}
