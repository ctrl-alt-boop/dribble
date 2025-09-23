package postgres

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

type SSLMode string

const (
	SSLModeDisable  SSLMode = "disable"
	SSLModeRequire  SSLMode = "require"
	SSLModeVerify   SSLMode = "verify-ca"
	SSLModeVerifyCA SSLMode = "verify-full"
)

type PostgreSQLDSN struct {
	Addr     string  `json:"addr"`
	Port     int     `json:"port"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	DBName   string  `json:"dbname"`
	SSLMode  SSLMode `json:"sslmode"`
}

// Type implements database.DataSourceNamer.
func (p PostgreSQLDSN) Type() database.Type {
	return database.PostgreSQL
}

func (p PostgreSQLDSN) DSN() string {
	// PostreSQL connection string format: host=localhost port=5432 user=postgres password=password dbname=mydb sslmode=disable

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		p.Addr, p.Port, p.Username, p.Password, p.DBName)
	if p.SSLMode != "" {
		dsn += " sslmode=" + string(p.SSLMode)
	} else {
		dsn += " sslmode=disable"
	}
	return dsn
}

// DSNOption defines a function that configures a PostgreSQLDSN.
type DSNOption func(*PostgreSQLDSN)

// NewDSN creates a new PostgreSQLDSN with the given options.
func NewDSN(opts ...DSNOption) *PostgreSQLDSN {
	dsn := &PostgreSQLDSN{
		Addr:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "",
		SSLMode:  SSLModeDisable,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// WithAddr sets the address for the DSN.
func WithAddr(addr string) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.Addr = addr
	}
}

// WithPort sets the port for the DSN.
func WithPort(port int) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.Port = port
	}
}

// WithUsername sets the username for the DSN.
func WithUsername(username string) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.Username = username
	}
}

// WithPassword sets the password for the DSN.
func WithPassword(password string) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.Password = password
	}
}

// WithDBName sets the database name for the DSN.
func WithDBName(dbname string) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.DBName = dbname
	}
}

// WithSSLMode sets the SSL mode for the DSN.
func WithSSLMode(sslmode SSLMode) DSNOption {
	return func(p *PostgreSQLDSN) {
		p.SSLMode = sslmode
	}
}
