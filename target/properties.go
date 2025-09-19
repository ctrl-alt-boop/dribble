package target

import (
	"maps"

	"github.com/ctrl-alt-boop/dribble/database"
)

type Properties struct {
	Dialect    database.SQLDialect
	Ip         string
	Port       int
	DBName     string
	Username   string
	Password   string // FIXME: some other type
	Additional map[string]string
}

func (s *Target) Copy(opts ...Option) *Target {
	newTarget := &Target{
		Name:       s.Name,
		Type:       s.Type,
		Properties: s.Properties,
	}

	maps.Copy(newTarget.Properties.Additional, s.Properties.Additional)

	for _, opt := range opts {
		opt(newTarget)
	}

	return newTarget
}

type Option func(*Target)

func AsTableSelect(table string) Option {
	return func(target *Target) {
		target.Type = TableTable
		target.Properties.Additional["select"] = "SELECT * FROM " + table
	}
}

func AsQuery(query string) Option {
	return func(target *Target) {
		target.Type = TableTable
		target.Properties.Additional["query"] = query
	}
}

func WithIp(ip string) Option {
	return func(target *Target) {
		target.Properties.Ip = ip
	}
}

func WithPort(port int) Option {
	return func(target *Target) {
		target.Properties.Port = port
	}
}

func WithHost(hostname string, port int) Option {
	return func(target *Target) {
		target.Properties.Ip = hostname
		target.Properties.Port = port
	}
}

func WithDB(name string) Option {
	return func(target *Target) {
		target.Properties.DBName = name
	}
}

func WithUser(user string) Option {
	return func(target *Target) {
		target.Properties.Username = user
	}
}

func WithPassword(pass string) Option {
	return func(target *Target) {
		target.Properties.Password = pass
	}
}

func WithProperty(key, value string) Option {
	return func(target *Target) {
		target.Properties.Additional[key] = value
	}
}
