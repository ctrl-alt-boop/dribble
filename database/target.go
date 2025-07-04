package database

import (
	"maps"
)

type TargetType string

const (
	None     TargetType = ""
	DBDriver TargetType = "Driver"
	DBServer TargetType = "Server"
	Database TargetType = "Database"
	DBTable  TargetType = "Table"
	Unknown  TargetType = "???"
)

type Target struct {
	Name               string
	Type               TargetType
	DriverName         string
	Ip                 string
	Port               int
	DBName             string
	Username           string
	Password           string // FIXME: some other type?
	AdditionalSettings map[string]string
}

func (s Target) Copy(opts ...TargetOption) *Target {
	newTarget := &Target{
		Name:               s.Name,
		Type:               s.Type,
		DriverName:         s.DriverName,
		Ip:                 s.Ip,
		Port:               s.Port,
		DBName:             s.DBName,
		Username:           s.Username,
		Password:           s.Password,
		AdditionalSettings: make(map[string]string),
	}

	maps.Copy(newTarget.AdditionalSettings, s.AdditionalSettings)

	for _, opt := range opts {
		opt(newTarget)
	}

	return newTarget
}

func NewTarget(name string, options ...TargetOption) *Target {
	settings := &Target{
		Name:               name,
		Type:               None,
		DriverName:         "",
		Ip:                 "localhost",
		Port:               0,
		DBName:             "",
		Username:           "",
		Password:           "",
		AdditionalSettings: make(map[string]string),
	}

	for _, option := range options {
		option(settings)
	}

	return settings
}

type TargetOption func(*Target)

func AsType(connectionType TargetType) TargetOption {
	return func(target *Target) {
		target.Type = connectionType
	}
}

func AsTableSelect(table string) TargetOption {
	return func(target *Target) {
		target.Type = DBTable
		target.AdditionalSettings["select"] = "SELECT * FROM " + table
	}
}

func WithDriver(name string) TargetOption {
	return func(target *Target) {
		target.DriverName = name
	}
}

func WithIp(ip string) TargetOption {
	return func(target *Target) {
		target.Ip = ip
	}
}

func WithPort(port int) TargetOption {
	return func(target *Target) {
		target.Port = port
	}
}

func WithHost(hostname string, port int) TargetOption {
	return func(target *Target) {
		target.Ip = hostname
		target.Port = port
	}
}

func WithDB(name string) TargetOption {
	return func(target *Target) {
		target.DBName = name
	}
}

func WithUser(user string) TargetOption {
	return func(target *Target) {
		target.Username = user
	}
}

func WithPassword(pass string) TargetOption {
	return func(target *Target) {
		target.Password = pass
	}
}

func WithSetting(key, value string) TargetOption {
	return func(target *Target) {
		target.AdditionalSettings[key] = value
	}
}
