package explorer

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/core/util"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

const commandlineID string = string(ExplorerPageID) + ".commandline"

// CommandlineMsg is sent from commandline when submitting a command
type CommandlineMsg struct {
	value string
}

type commandline struct {
	input textinput.Model

	width int

	showSuggestions bool

	style lipgloss.Style

	keybinds *CommandlineKeys
}

func newCommandline() *commandline {
	return &commandline{
		input:           textinput.New(),
		showSuggestions: false,

		keybinds: DefaultCommandlineKeyBindings(),
	}
}

func (p *commandline) SetWidth(width int) {
	p.width = width
}

func (p *commandline) ShowSuggestions() {
	p.showSuggestions = true
}

func (p *commandline) Init() tea.Cmd {
	p.input = textinput.New()
	p.input.EchoMode = textinput.EchoNormal
	p.input.Prompt = "> "
	p.input.CharLimit = p.width - 1

	p.input.SetSuggestions(util.GetSqlSuggestions()) // FIXME: TEMP

	return nil
}

func (p *commandline) Update(msg tea.Msg) (*commandline, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case focusChangedMsg:
		if string(msg) == commandlineID {
			return p, p.input.Focus()
		} else {
			p.input.Reset()
			p.input.Blur()
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, p.keybinds.Cancel):
			return p, focusBackwards

		case key.Matches(msg, p.keybinds.Submit):
			logging.GlobalLogger().Infof("cmd> %s", p.input.Value())
			return p, tea.Batch(focusBackwards, p.submitCommand)

		}
	}
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p *commandline) submitCommand() tea.Msg {
	return &CommandlineMsg{
		value: p.input.Value(),
	}
}

func (p *commandline) SetStyle(style lipgloss.Style) {
	p.style = style.UnsetPadding().UnsetMargins()
}

func (p *commandline) Render() *lipgloss.Layer {
	p.input.SetWidth(p.width)
	style := p.style.
		Width(p.width).Height(3)
	return lipgloss.NewLayer(style.Render(p.input.View()))
}
