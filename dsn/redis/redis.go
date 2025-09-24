package redis

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.DataSourceNamer = (*RedisDSN)(nil)

type RedisDSN struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"` // Redis DB number
}

// Info implements database.DataSourceNamer.
func (r *RedisDSN) Info() string {
	if r.DB == 0 {
		return fmt.Sprintf("Redis: %s:%d", r.Addr, r.Port)
	}
	return fmt.Sprintf("Redis: %s:%d/%d", r.Addr, r.Port, r.DB)
}

// Type implements database.DataSourceNamer.
func (r RedisDSN) Type() database.Type {
	return database.Redis
}

func (r RedisDSN) DSN() string {
	// Redis connection string format: redis://[username:password@]host:port[/db]

	var parts []string
	if r.Username != "" {
		parts = append(parts, r.Username)
		if r.Password != "" {
			parts = append(parts, ":", r.Password)
		}
		parts = append(parts, "@")
	}
	parts = append(parts, r.Addr)
	if r.Port != 0 {
		parts = append(parts, fmt.Sprintf(":%d", r.Port))
	}
	if r.DB != 0 {
		parts = append(parts, fmt.Sprintf("/%d", r.DB))
	}

	return "redis://" + strings.Join(parts, "")
}

// DSNOption defines a function that configures a RedisDSN.
type DSNOption func(*RedisDSN)

// NewDSN creates a new RedisDSN with the given options.
func NewDSN(opts ...DSNOption) *RedisDSN {
	dsn := &RedisDSN{
		Addr: "localhost",
		Port: 6379,
		DB:   0,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// WithAddr sets the address for the DSN.
func WithAddr(addr string) DSNOption {
	return func(r *RedisDSN) {
		r.Addr = addr
	}
}

// WithPort sets the port for the DSN.
func WithPort(port int) DSNOption {
	return func(r *RedisDSN) {
		r.Port = port
	}
}

// WithUsername sets the username for the DSN.
func WithUsername(username string) DSNOption {
	return func(r *RedisDSN) {
		r.Username = username
	}
}

// WithPassword sets the password for the DSN.
func WithPassword(password string) DSNOption {
	return func(r *RedisDSN) {
		r.Password = password
	}
}

// WithDB sets the database number for the DSN.
func WithDB(db int) DSNOption {
	return func(r *RedisDSN) {
		r.DB = db
	}
}
