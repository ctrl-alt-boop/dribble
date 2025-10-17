package component

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui"
	"github.com/ctrl-alt-boop/dribbler/ui/content"
)

type Panel struct {
	list content.List

	Width, Height int

	showDetails bool

	isLoading bool
	spinner   spinner.Model

	selectIndexHistory []int
}

func NewPanel() *Panel {
	return &Panel{
		list:               content.NewList([]content.Item{}),
		spinner:            spinner.New(spinner.WithSpinner(ui.MovingBlock)),
		selectIndexHistory: make([]int, 0),
		showDetails:        config.Cfg.Ui.ShowTargetDetails,
	}
}

func (p *Panel) GetSelected() any {
	return p.list.GetSelected()
}

func (p Panel) Init() tea.Cmd {
	return nil
}

func (p Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := p
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updated.Width, updated.Height = msg.Width, msg.Height
		return updated, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Up):
			updated.list.MoveCursorUp()

		case key.Matches(msg, config.Keys.Down):
			updated.list.MoveCursorDown()

		case key.Matches(msg, config.Keys.Select):
			// return p, p.OnSelect()

		case key.Matches(msg, config.Keys.Back):
			return updated, nil

		case key.Matches(msg, config.Keys.Details):
			updated.showDetails = !p.showDetails
			// return p, p.OnShowTableDetails()

		case key.Matches(msg, config.Keys.New):

		case key.Matches(msg, config.Keys.NewEmpty):
			return updated, func() tea.Msg {
				return OpenIntentBuilderMsg{Method: database.Read, Table: ""}
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		updated.spinner, cmd = p.spinner.Update(msg)
		return updated, cmd
	}

	return updated, nil
}

func (p Panel) View() string {
	if p.isLoading {
		return lipgloss.Place(
			p.Width,
			p.Height,
			lipgloss.Center,
			lipgloss.Position(0.2),
			p.spinner.View(),
		)
	}

	// details, ok := p.list.GetSelected().(*ui.ConnectionItem)
	// var detailsView string
	// if p.showDetails && ok {
	// 	detailsView = lipgloss.NewStyle().
	// 		BorderStyle(lipgloss.NormalBorder()).BorderTop(true).
	// 		Padding(1, 0).
	// 		MaxHeight(9).Width(p.Width).
	// 		Render(details.Inspect())
	// }

	return p.list.View()
}
