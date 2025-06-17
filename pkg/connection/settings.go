package connection

import (
	"bytes"
	"maps"
	"text/template"
)

type Type string

const (
	TypeUnknown Type = "(none)"
	Driver      Type = "Driver"
	Server      Type = "Server"
	Database    Type = "Database"
	Table       Type = "Table"
)

type Settings struct {
	SettingsName       string
	Type               Type
	DriverName         string
	Ip                 string
	Port               int
	DbName             string
	Username           string
	Password           string
	AdditionalSettings map[string]string
}

func (s Settings) Copy(opts ...Option) *Settings {
	newSettings := &Settings{
		SettingsName:       s.SettingsName,
		Type:               s.Type,
		DriverName:         s.DriverName,
		Ip:                 s.Ip,
		Port:               s.Port,
		DbName:             s.DbName,
		Username:           s.Username,
		Password:           s.Password,
		AdditionalSettings: make(map[string]string),
	}

	maps.Copy(newSettings.AdditionalSettings, s.AdditionalSettings)

	for _, opt := range opts {
		opt(newSettings)
	}

	return newSettings
}

func NewSettings(options ...Option) *Settings {
	settings := &Settings{
		SettingsName:       "",
		Type:               TypeUnknown,
		DriverName:         "",
		Ip:                 "localhost",
		Port:               0,
		DbName:             "",
		Username:           "",
		Password:           "",
		AdditionalSettings: make(map[string]string),
	}

	for _, option := range options {
		option(settings)
	}

	return settings
}

const stringTemplate = `{{.DriverName}}
{{.Username}}:********
{{.Ip}}:{{.Port}}
{{- if .DbName}}
{{.DbName}}
{{- end -}}`

func (s Settings) AsString() string {
	var buf bytes.Buffer
	template.Must(template.New("settings").Parse(stringTemplate)).Execute(&buf, s)
	return buf.String()
}

type Option func(*Settings)

func AsType(connectionType Type) Option {
	return func(settings *Settings) {
		settings.Type = connectionType
	}
}

func AsTableSelect(table string) Option {
	return func(settings *Settings) {
		settings.Type = Table
		settings.AdditionalSettings["select"] = "SELECT * FROM " + table
	}
}

func WithDriver(name string) Option {
	return func(settings *Settings) {
		settings.DriverName = name
	}
}

func WithIp(ip string) Option {
	return func(settings *Settings) {
		settings.Ip = ip
	}
}

func WithPort(port int) Option {
	return func(settings *Settings) {
		settings.Port = port
	}
}

func WithHost(hostname string, port int) Option {
	return func(settings *Settings) {
		settings.Ip = hostname
		settings.Port = port
	}
}

func WithDb(name string) Option {
	return func(settings *Settings) {
		settings.DbName = name
	}
}

func WithUser(user string) Option {
	return func(settings *Settings) {
		settings.Username = user
	}
}

func WithPassword(pass string) Option {
	return func(settings *Settings) {
		settings.Password = pass
	}
}

func WithSetting(key, value string) Option {
	return func(settings *Settings) {
		settings.AdditionalSettings[key] = value
	}
}
