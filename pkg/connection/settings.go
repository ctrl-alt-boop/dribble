package connection

import (
	"bytes"
	"maps"
	"text/template"
)

type Settings struct {
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
{{.DbName}}
`

func (s Settings) AsString() string {
	var buf bytes.Buffer
	template.Must(template.New("settings").Parse(stringTemplate)).Execute(&buf, s)
	return buf.String()
}

type Option func(*Settings)

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
