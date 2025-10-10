package dribbler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

type Model struct {
	Width, Height           int
	InnerWidth, InnerHeight int

	MainContent widget.ContentArea

	dribbleClient *dribble.Client
}

func NewDribblerModel() Model {
	return Model{
		MainContent: CreateDemoLayout(),
	}
}

func (m *Model) Init() tea.Cmd {
	config.LoadConfig()
	m.dribbleClient = dribble.NewClient()

	return m.MainContent.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.MainContent.View()
}
