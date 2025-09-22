package database

import (
	"maps"
)

type ConnectionProperties struct {
	Addr       string
	Port       int
	DBName     string
	Username   string
	Password   string // FIXME: some other type
	Additional map[string]string
}

func NewConnectionProperties(opts ...Option) *ConnectionProperties {
	properties := &ConnectionProperties{
		Addr:       "localhost",
		Port:       0,
		DBName:     "",
		Username:   "",
		Password:   "",
		Additional: make(map[string]string),
	}

	for _, opt := range opts {
		opt(properties)
	}

	return properties
}

func (p *ConnectionProperties) Copy(opts ...Option) *ConnectionProperties {
	newProperties := &ConnectionProperties{
		Addr:       p.Addr,
		Port:       p.Port,
		DBName:     p.DBName,
		Username:   p.Username,
		Password:   p.Password,
		Additional: make(map[string]string),
	}

	maps.Copy(newProperties.Additional, p.Additional)

	for _, opt := range opts {
		opt(newProperties)
	}

	return newProperties
}

type Option func(*ConnectionProperties) // FIXME: Rename please

func WithIp(ip string) Option {
	return func(prop *ConnectionProperties) {
		prop.Addr = ip
	}
}

func WithPort(port int) Option {
	return func(prop *ConnectionProperties) {
		prop.Port = port
	}
}

func WithHost(hostname string, port int) Option {
	return func(prop *ConnectionProperties) {
		prop.Addr = hostname
		prop.Port = port
	}
}

func WithDB(name string) Option {
	return func(prop *ConnectionProperties) {
		prop.DBName = name
	}
}

func WithUser(user string) Option {
	return func(prop *ConnectionProperties) {
		prop.Username = user
	}
}

func WithPassword(pass string) Option {
	return func(prop *ConnectionProperties) {
		prop.Password = pass
	}
}

func WithProperty(key, value string) Option {
	return func(prop *ConnectionProperties) {
		prop.Additional[key] = value
	}
}
