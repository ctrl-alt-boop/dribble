package dsn

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql/postgres"
)

var _ datasource.Namer = (*PostgreSQL)(nil)

type SSLMode string

const (
	SSLModeDisable  SSLMode = "disable"
	SSLModeRequire  SSLMode = "require"
	SSLModeVerify   SSLMode = "verify-ca"
	SSLModeVerifyCA SSLMode = "verify-full"
)

type PostgreSQL struct {
	Addr     string  `json:"addr"`
	Port     int     `json:"port"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	DBName   string  `json:"dbname"`
	SSLMode  SSLMode `json:"sslmode"`
}

// SourceType implements datasource.Namer.
func (p *PostgreSQL) SourceType() datasource.SourceType {
	return postgres.SourceType
}

// Info implements database.DataSourceNamer.
func (p *PostgreSQL) Info() string {
	if p.DBName == "" {
		return fmt.Sprintf("PostgreSQL: %s:%d", p.Addr, p.Port)
	}
	return fmt.Sprintf("PostgreSQL: %s:%d/%s", p.Addr, p.Port, p.DBName)
}

// Type implements database.DataSourceNamer.
func (p PostgreSQL) Type() datasource.Type {
	return datasource.PostgreSQL
}

func (p PostgreSQL) DSN() string {
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

// PostgresOption defines a function that configures a PostgreSQLDSN.
type PostgresOption func(*PostgreSQL)

// PostgresDSN creates a new PostgreSQLDSN with the given options.
func PostgresDSN(opts ...PostgresOption) *PostgreSQL {
	dsn := &PostgreSQL{
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

// PostgresAddr sets the address for the DSN.
func PostgresAddr(addr string) PostgresOption {
	return func(p *PostgreSQL) {
		p.Addr = addr
	}
}

// PostgresPort sets the port for the DSN.
func PostgresPort(port int) PostgresOption {
	return func(p *PostgreSQL) {
		p.Port = port
	}
}

// PostgresUsername sets the username for the DSN.
func PostgresUsername(username string) PostgresOption {
	return func(p *PostgreSQL) {
		p.Username = username
	}
}

// PostgresPassword sets the password for the DSN.
func PostgresPassword(password string) PostgresOption {
	return func(p *PostgreSQL) {
		p.Password = password
	}
}

// PostgresDBName sets the database name for the DSN.
func PostgresDBName(dbname string) PostgresOption {
	return func(p *PostgreSQL) {
		p.DBName = dbname
	}
}

// PostgresSSLMode sets the SSL mode for the DSN.
func PostgresSSLMode(sslmode SSLMode) PostgresOption {
	return func(p *PostgreSQL) {
		p.SSLMode = sslmode
	}
}
