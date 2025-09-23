package mysql

// import (
// 	"fmt"

// 	"github.com/ctrl-alt-boop/dribble/database"
// )

// type MySQLDSN struct {
// 	Addr     string `json:"addr"`
// 	Port     int    `json:"port"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	DBName   string `json:"dbname"`
// }

// // Type implements database.DataSourceNamer.
// func (m MySQLDSN) Type() database.Type {
// 	return database.MySQL
// }

// func (m MySQLDSN) DSN() string {
// 	dsn := ""
// 	if m.Username != "" {
// 		dsn += m.Username
// 		if m.Password != "" {
// 			dsn += ":" + m.Password
// 		}
// 		dsn += "@"
// 	}
// 	dsn += "tcp(" + m.Addr
// 	if m.Port != 0 {
// 		dsn += fmt.Sprintf(":%d", m.Port)
// 	}
// 	dsn += ")"
// 	if m.DBName != "" {
// 		dsn += "/" + m.DBName
// 	}
// 	return dsn
// }

// // DSNOption defines a function that configures a MySQLDSN.
// type DSNOption func(*MySQLDSN)

// // NewMySQLDSN creates a new MySQLDSN with the given options.
// func NewMySQLDSN(opts ...DSNOption) *MySQLDSN {
// 	dsn := &MySQLDSN{
// 		Addr:     "localhost",
// 		Port:     3306,
// 		Username: "root",
// 		Password: "",
// 	}
// 	for _, opt := range opts {
// 		opt(dsn)
// 	}
// 	return dsn
// }

// // WithAddr sets the address for the DSN.
// func WithAddr(addr string) DSNOption {
// 	return func(m *MySQLDSN) {
// 		m.Addr = addr
// 	}
// }

// // WithPort sets the port for the DSN.
// func WithPort(port int) DSNOption {
// 	return func(m *MySQLDSN) {
// 		m.Port = port
// 	}
// }

// // WithUsername sets the username for the DSN.
// func WithUsername(username string) DSNOption {
// 	return func(m *MySQLDSN) {
// 		m.Username = username
// 	}
// }

// // WithPassword sets the password for the DSN.
// func WithPassword(password string) DSNOption {
// 	return func(m *MySQLDSN) {
// 		m.Password = password
// 	}
// }

// // WithDBName sets the database name for the DSN.
// func WithDBName(dbname string) DSNOption {
// 	return func(m *MySQLDSN) {
// 		m.DBName = dbname
// 	}
// }
