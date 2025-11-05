package dribbler

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
	"github.com/ctrl-alt-boop/dribbler/internal/page/explorer"
	"github.com/ctrl-alt-boop/dribbler/keys"
)

type Model struct {
	Width, Height           int
	InnerWidth, InnerHeight int

	currentPage page.Page
}

func NewModel() Model {
	return Model{
		currentPage: explorer.NewExplorerPage(250, 50),
	}
}

func (m Model) Init() tea.Cmd {
	// load config then check if there is a config.startup.target or similar

	//

	return m.currentPage.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg: // Handle application wide key presses
		if key.Matches(msg, keys.Map.Quit) {
			return m, tea.Quit
		}
	}

	updatedPage, cmd := m.currentPage.Update(msg)
	m.currentPage = updatedPage

	return m, cmd
}

func (m Model) View() tea.View {
	canvas := m.currentPage.Render()

	view := tea.NewView(canvas.Render())
	view.WindowTitle = "Dribbler"
	view.AltScreen = true

	return view
}
