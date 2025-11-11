package dsn

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql/mysql"
)

var _ datasource.Namer = (*MySQL)(nil)

type MySQL struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// SourceType implements datasource.Namer.
func (m *MySQL) SourceType() datasource.SourceType {
	return mysql.SourceType
}

// Info implements database.DataSourceNamer.
func (m *MySQL) Info() string {
	if m.DBName == "" {
		return fmt.Sprintf("MySQL: %s:%d", m.Addr, m.Port)
	}
	return fmt.Sprintf("MySQL: %s:%d/%s", m.Addr, m.Port, m.DBName)
}

// Type implements database.DataSourceNamer.
func (m MySQL) Type() datasource.Type {
	return datasource.MySQL
}

func (m MySQL) DSN() string {
	dsn := ""
	if m.Username != "" {
		dsn += m.Username
		if m.Password != "" {
			dsn += ":" + m.Password
		}
		dsn += "@"
	}
	dsn += "tcp(" + m.Addr
	if m.Port != 0 {
		dsn += fmt.Sprintf(":%d", m.Port)
	}
	dsn += ")" + "/"
	if m.DBName != "" {
		dsn += m.DBName
	}
	return dsn
}

// MySQLOption defines a function that configures a MySQL DSN.
type MySQLOption func(*MySQL)

// MySQLDSN creates a new MySQLDSN with the given options.
func MySQLDSN(opts ...MySQLOption) *MySQL {
	dsn := &MySQL{
		Addr:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "",
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// MySQLAddr sets the address for the DSN.
func MySQLAddr(addr string) MySQLOption {
	return func(m *MySQL) {
		m.Addr = addr
	}
}

// MySQLPort sets the port for the DSN.
func MySQLPort(port int) MySQLOption {
	return func(m *MySQL) {
		m.Port = port
	}
}

// MySQLUsername sets the username for the DSN.
func MySQLUsername(username string) MySQLOption {
	return func(m *MySQL) {
		m.Username = username
	}
}

// MySQLPassword sets the password for the DSN.
func MySQLPassword(password string) MySQLOption {
	return func(m *MySQL) {
		m.Password = password
	}
}

// MySQLDBName sets the database name for the DSN.
func MySQLDBName(dbname string) MySQLOption {
	return func(m *MySQL) {
		m.DBName = dbname
	}
}
