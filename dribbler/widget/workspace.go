package widget

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble"

	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/io"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

type Workspace struct {
	viewport                    viewport.Model
	dribbleClient               *dribble.Client
	Width, Height               int
	ContentWidth, ContentHeight int

	table *ui.DribbleTable

	isLoading bool
	spinner   spinner.Model
}

func NewWorkspace(dribbleClient *dribble.Client) *Workspace {
	return &Workspace{
		dribbleClient: dribbleClient,
		viewport:      viewport.New(0, 0),
		table:         ui.NewDribbleTable(),
	}
}

func (d *Workspace) Init() tea.Cmd {
	d.spinner = spinner.New()
	d.spinner.Spinner = ui.MovingBlock
	d.isLoading = false

	return nil
}

func (d *Workspace) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Up):
			d.table.MoveCursorUp()
		case key.Matches(msg, config.Keys.Down):
			d.table.MoveCursorDown()
		case key.Matches(msg, config.Keys.Left):
			d.table.MoveCursorLeft()
		case key.Matches(msg, config.Keys.Right):
			d.table.MoveCursorRight()
		case key.Matches(msg, config.Keys.Select):
			return d, d.SelectCell
		case key.Matches(msg, config.Keys.Increase):
			d.table.IncreaseColumnSize() // FIXME: Resizing the workspace
		case key.Matches(msg, config.Keys.Decrease):
			d.table.DecreaseColumnSize() // FIXME: Resizing the workspace
		}

	case io.DribbleEventMsg:
		d.isLoading = false
		switch msg.Type {
		case dribble.SuccessReadTable:
			args, ok := msg.Args.(dribble.TableFetchData)
			if ok {
				return d, tea.Batch(d.SetTable(*args.Table), RequestFocusChange(KindWorkspace))
			}
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
	return WorkspaceSet
}

func (d *Workspace) IsTableSet() bool {
	return d.table.IsTableSet()
}

func (d *Workspace) UpdateSize(width, height int) {
	d.Width, d.Height = width, height
	d.ContentWidth = d.Width - ui.WorkspaceStyle.GetHorizontalFrameSize()
	d.ContentHeight = d.Height - ui.WorkspaceStyle.GetVerticalFrameSize()
	d.viewport.Width = d.ContentWidth
	d.viewport.Height = d.ContentHeight

	d.table.ViewportWidth = d.ContentWidth
	d.table.ViewportHeight = d.ContentHeight
}

func (d *Workspace) ViewportWidth() int {
	return d.table.Width
}

func (d *Workspace) View() string {
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
