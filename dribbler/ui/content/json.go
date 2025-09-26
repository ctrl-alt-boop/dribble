package content

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
)

var _ Content[string] = (*Json)(nil)

type Json string

func (j Json) Data() any {
	return string(j)
}

func (j Json) Get() string {
	return j.PrettifyJson(string(j))
}

func (j *Json) UpdateSize(width int, height int) {}

func (j *Json) PrettifyJson(jsonString string) string {
	var tmp any
	if err := json.Unmarshal([]byte(jsonString), &tmp); err != nil {
		return jsonString
	}
	prettyJson, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return jsonString
	}
	return string(prettyJson)
}

func (j *Json) View() string {
	return j.Get()
}

func (j *Json) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return j, nil
}

func (j *Json) Init() tea.Cmd {
	return nil
}
