package dsn

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
)

var _ datasource.Namer = (*MongoDB)(nil)

type MongoDB struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// SourceType implements datasource.Namer.
func (m *MongoDB) SourceType() datasource.SourceType {
	panic("unimplemented")
}

// Info implements database.DataSourceNamer.
func (m MongoDB) Info() string {
	if m.DBName == "" {
		return fmt.Sprintf("MongoDB: %s:%d", m.Addr, m.Port)
	}
	return fmt.Sprintf("MongoDB: %s:%d/%s", m.Addr, m.Port, m.DBName)
}

// Type implements database.DataSourceNamer.
func (m MongoDB) Type() datasource.Type {
	return datasource.MongoDB
}

func (m MongoDB) DSN() string {
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

// MongoDBOption defines a function that configures a MongoDB DSN.
type MongoDBOption func(*MongoDB)

// MongoDSN creates a new MongoDBDSN with the given options.
func MongoDSN(opts ...MongoDBOption) *MongoDB {
	dsn := &MongoDB{
		Addr: "localhost",
		Port: 27017,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// MongoAddr sets the address for the DSN.
func MongoAddr(addr string) MongoDBOption {
	return func(m *MongoDB) {
		m.Addr = addr
	}
}

// MongoPort sets the port for the DSN.
func MongoPort(port int) MongoDBOption {
	return func(m *MongoDB) {
		m.Port = port
	}
}

// MongoUsername sets the username for the DSN.
func MongoUsername(username string) MongoDBOption {
	return func(m *MongoDB) {
		m.Username = username
	}
}

// MongoPassword sets the password for the DSN.
func MongoPassword(password string) MongoDBOption {
	return func(m *MongoDB) {
		m.Password = password
	}
}

// MongoDBName sets the database name for the DSN.
func MongoDBName(dbname string) MongoDBOption {
	return func(m *MongoDB) {
		m.DBName = dbname
	}
}
