package mongo

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

type MongoDBDSN struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Type implements database.DataSourceNamer.
func (m MongoDBDSN) Type() database.Type {
	return database.MongoDB
}

func (m MongoDBDSN) DSN() string {
	// MongoDB connection string format: mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]

	dsn := "mongodb://"
	if m.Username != "" {
		dsn += m.Username
		if m.Password != "" {
			dsn += ":" + m.Password
		}
		dsn += "@"
	}
	dsn += m.Addr
	if m.Port != 0 {
		dsn += fmt.Sprintf(":%d", m.Port)
	}
	if m.DBName != "" {
		dsn += "/" + m.DBName
	}
	return dsn
}

// DSNOption defines a function that configures a MongoDBDSN.
type DSNOption func(*MongoDBDSN)

// NewDSN creates a new MongoDBDSN with the given options.
func NewDSN(opts ...DSNOption) *MongoDBDSN {
	dsn := &MongoDBDSN{
		Addr: "localhost",
		Port: 27017,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// WithAddr sets the address for the DSN.
func WithAddr(addr string) DSNOption {
	return func(m *MongoDBDSN) {
		m.Addr = addr
	}
}

// WithPort sets the port for the DSN.
func WithPort(port int) DSNOption {
	return func(m *MongoDBDSN) {
		m.Port = port
	}
}

// WithUsername sets the username for the DSN.
func WithUsername(username string) DSNOption {
	return func(m *MongoDBDSN) {
		m.Username = username
	}
}

// WithPassword sets the password for the DSN.
func WithPassword(password string) DSNOption {
	return func(m *MongoDBDSN) {
		m.Password = password
	}
}

// WithDBName sets the database name for the DSN.
func WithDBName(dbname string) DSNOption {
	return func(m *MongoDBDSN) {
		m.DBName = dbname
	}
}
