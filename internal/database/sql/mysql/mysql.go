package mysql

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	_ "github.com/go-sql-driver/mysql"
)

var _ database.Driver = &MySQL{}

type MySQL struct{}

func NewMySQLDriver(target *database.Target) (*MySQL, error) {
	driver := &MySQL{}
	return driver, nil
}

// // Capabilities implements database.Dialect.
func (m *MySQL) Capabilities() []database.Capabilities {
	return []database.Capabilities{}
}

// ConnectionString implements database.Driver.
// Server=myServerAddress;Port=1234;Database=myDataBase;Uid=myUsername;Pwd=myPassword;
func (m *MySQL) ConnectionString(target *database.Target) string {
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
func (m *MySQL) Dialect() database.Dialect {
	return m
}

// RenderIntent implements database.Driver.
func (m *MySQL) RenderIntent(intent *database.Intent) (string, error) {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (m *MySQL) ResolveType(dbType string, value []byte) (any, error) {
	return string(value), nil
}
