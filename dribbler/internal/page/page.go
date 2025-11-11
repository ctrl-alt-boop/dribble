package page

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type ID string

type Page interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Page, tea.Cmd)
	Render() *lipgloss.Canvas

	SetSize(width, height int)
}
