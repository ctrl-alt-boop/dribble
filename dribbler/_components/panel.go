package components

import (
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/keys"
)

var (
	movingBlockFrames = []string{
		"██        ", " ██       ", "  ██      ", "   ██     ", "    ██    ", "     ██   ", "      ██  ", "       ██ ",
		"        ██", "       ██ ", "      ██  ", "     ██   ", "    ██    ", "   ██     ", "  ██      ", " ██       ",
	}
	MovingBlock = spinner.Spinner{
		Frames: movingBlockFrames,
		FPS:    time.Second / 18, //nolint:mnd
	}

	GrowingBlock = spinner.Spinner{
		Frames: []string{"█", "███", "█████", "███████", "█████████", "███████", "█████", "███", "█"},
		FPS:    time.Second / 10, //nolint:mnd
	}
)

type Panel struct {
	// content T

	Width, Height int

	showDetails bool

	isLoading bool
	spinner   spinner.Model

	selectIndexHistory []int
}

func NewPanel() *Panel {
	return &Panel{
		spinner:            spinner.New(spinner.WithSpinner(MovingBlock)),
		selectIndexHistory: make([]int, 0),
		showDetails:        config.Cfg.Ui.ShowTargetDetails,
	}
}

func (p *Panel) GetSelected() any {
	return nil
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
		case key.Matches(msg, keys.Map.Up):
			// updated.list.MoveCursorUp()

		case key.Matches(msg, keys.Map.Down):
			// updated.list.MoveCursorDown()

		case key.Matches(msg, keys.Map.Select):
			// return p, p.OnSelect()

		case key.Matches(msg, keys.Map.Back):
			return updated, nil

		case key.Matches(msg, keys.Map.Details):
			updated.showDetails = !p.showDetails
			// return p, p.OnShowTableDetails()

		case key.Matches(msg, keys.Map.New):

		case key.Matches(msg, keys.Map.NewEmpty):
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

func (p Panel) Render() string {
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

	return "" // p.list.Render()
}

func (p Panel) View() tea.View {
	return tea.NewView(p.Render())
}
