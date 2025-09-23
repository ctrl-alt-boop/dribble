package ui

import (
	"html/template"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
)

const (
	PanelWidthRatio       int = 6
	BorderThickness       int = 1
	BorderThicknessDouble int = 2
	PopupMargin           int = 10
)

const stringTemplate = `{{.DriverName}}
{{.Username}}:********
{{.Ip}}:{{.Port}}
{{- if .DBName}}
{{.DBName}}
{{- end -}}`

func AsString(t database.DataSourceNamer) string {
	var sb strings.Builder
	template.Must(template.New("settings").Parse(stringTemplate)).Execute(&sb, t)
	return sb.String()
}
