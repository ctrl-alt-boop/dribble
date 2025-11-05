package explorer

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/core/util"
)

const promptPanelID string = string(ExplorerPageID) + ".prompt"

var esc = key.NewBinding(
	key.WithKeys("esc"),
)

type PromptMsg struct {
	Value string
}

type promptPanel struct {
	input textinput.Model

	width int

	showSuggestions bool
}

func newPrompt() *promptPanel {
	return &promptPanel{
		input:           textinput.New(),
		showSuggestions: false,
	}
}

func (p *promptPanel) SetInnerWidth(width int) {
	p.width = width
}

func (p *promptPanel) ShowSuggestions() {
	p.showSuggestions = true
}

func (p *promptPanel) Init() tea.Cmd {
	p.input = textinput.New()
	p.input.EchoMode = textinput.EchoNormal
	p.input.Prompt = ">"
	p.input.CharLimit = p.width - 1

	p.input.SetSuggestions(util.GetSqlSuggestions()) // FIXME: TEMP

	return nil
}

func (p *promptPanel) Update(msg tea.Msg) (*promptPanel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, esc):
			p.input.Reset()
			p.input.Blur()
			return p, focusBackwards

		default:
			p.input, cmd = p.input.Update(msg)
			return p, cmd
		}
	}
	return p, nil
}

func (p *promptPanel) Render() *lipgloss.Layer {
	p.input.SetWidth(p.width)
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		Width(p.width).Height(3)
	return lipgloss.NewLayer(box.Render(p.input.View()))
}
