package dsn

import (
	"fmt"
	"strings"

	"github.com/ctrl-alt-boop/dribble/datasource"
)

var _ datasource.Namer = (*Redis)(nil)

type Redis struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"` // Redis DB number
}

// SourceType implements datasource.Namer.
func (r *Redis) SourceType() datasource.SourceType {
	panic("unimplemented")
}

// Info implements database.DataSourceNamer.
func (r *Redis) Info() string {
	if r.DB == 0 {
		return fmt.Sprintf("Redis: %s:%d", r.Addr, r.Port)
	}
	return fmt.Sprintf("Redis: %s:%d/%d", r.Addr, r.Port, r.DB)
}

// Type implements database.DataSourceNamer.
func (r Redis) Type() datasource.Type {
	return datasource.Redis
}

func (r Redis) DSN() string {
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

// RedisOption defines a function that configures a Redis DSN.
type RedisOption func(*Redis)

// RedisDSN creates a new RedisDSN with the given options.
func RedisDSN(opts ...RedisOption) *Redis {
	dsn := &Redis{
		Addr: "localhost",
		Port: 6379,
		DB:   0,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// RedisAddr sets the address for the DSN.
func RedisAddr(addr string) RedisOption {
	return func(r *Redis) {
		r.Addr = addr
	}
}

// RedisPort sets the port for the DSN.
func RedisPort(port int) RedisOption {
	return func(r *Redis) {
		r.Port = port
	}
}

// RedisUsername sets the username for the DSN.
func RedisUsername(username string) RedisOption {
	return func(r *Redis) {
		r.Username = username
	}
}

// RedisPassword sets the password for the DSN.
func RedisPassword(password string) RedisOption {
	return func(r *Redis) {
		r.Password = password
	}
}

// RedisDB sets the database number for the DSN.
func RedisDB(db int) RedisOption {
	return func(r *Redis) {
		r.DB = db
	}
}
