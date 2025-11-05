package page

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ID string

type Page interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Page, tea.Cmd)
	Render() *lipgloss.Canvas
}
