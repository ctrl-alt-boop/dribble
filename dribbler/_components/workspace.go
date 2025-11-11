package components

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/keys"

	"github.com/ctrl-alt-boop/dribble/result"
)

type Workspace struct {
	viewport                    viewport.Model
	dribbleClient               *dribble.Client
	Width, Height               int
	ContentWidth, ContentHeight int

	table *DribbleTable

	isLoading bool
	spinner   spinner.Model
}

func NewWorkspace() *Workspace {
	return &Workspace{
		viewport: viewport.New(viewport.WithWidth(0), viewport.WithHeight(0)),
		table:    NewDribbleTable(),
	}
}

func (d *Workspace) Init() tea.Cmd {
	d.spinner = spinner.New()
	d.spinner.Spinner = MovingBlock
	d.isLoading = false

	return nil
}

func (d *Workspace) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Map.Up):
			d.table.MoveCursorUp()
		case key.Matches(msg, keys.Map.Down):
			d.table.MoveCursorDown()
		case key.Matches(msg, keys.Map.Left):
			d.table.MoveCursorLeft()
		case key.Matches(msg, keys.Map.Right):
			d.table.MoveCursorRight()
		case key.Matches(msg, keys.Map.Select):
			return d, d.SelectCell
		case key.Matches(msg, keys.Map.Increase):
			d.table.IncreaseColumnSize() // FIXME: Resizing the workspace
		case key.Matches(msg, keys.Map.Decrease):
			d.table.DecreaseColumnSize() // FIXME: Resizing the workspace
		}
	}

	return d, nil
}

func (d *Workspace) SelectCell() tea.Msg {
	return OpenCellDataMsg{Value: d.table.GetSelected()}
}

func (d *Workspace) SetTable(table result.Table) tea.Cmd {
	d.table.SetTable(table)
	d.viewport.SetXOffset(0)
	d.viewport.SetYOffset(0)
	return nil
}

func (d *Workspace) Render() string {
	if d.isLoading {
		return lipgloss.Place(
			d.ContentWidth,
			d.ContentHeight,
			lipgloss.Center,
			lipgloss.Center,
			d.spinner.View(),
		)
	}

	d.viewport.SetContent(d.table.View())

	d.viewport.SetYOffset(d.table.Offset.Y)
	d.viewport.SetXOffset(d.table.Offset.X)

	return d.viewport.View()
	// return lipgloss.NewStyle().Width(d.ContentWidth).Height(d.ContentHeight).Render(d.viewport.View())
}

func (d *Workspace) View() tea.View {
	return tea.NewView(d.Render())
}
