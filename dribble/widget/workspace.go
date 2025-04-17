package widget

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/io"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
)

var (
	WorkspaceBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Right:       "│",
		Left:        "│",
		TopLeft:     "┬",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┤",
	}

	workspaceStyle = lipgloss.NewStyle().
			Border(WorkspaceBorder, true, true, false, true).
			Align(lipgloss.Left, lipgloss.Top)
)

type Workspace struct {
	width, height int
	goolDb        *gooldb.GoolDb

	table        *ui.BubblesTable
	columnWidths []int

	isLoading bool
	spinner   spinner.Model
}

func NewWorkspace(gool *gooldb.GoolDb) *Workspace {
	return &Workspace{
		goolDb:       gool,
		table:        ui.New(),
		columnWidths: make([]int, 0),
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
	case tea.WindowSizeMsg:
		d.UpdateSize(msg.Width, msg.Height)
	case io.GoolDbEventMsg:
		d.isLoading = false
		switch msg.Type {
		case gooldb.TableFetched:
			args, ok := msg.Args.(gooldb.TableFetchData)
			if ok {
				return d, d.SetTable(args.Table)
			}
		}
	}

	return d, nil
}

func (d *Workspace) SetTable(table data.Table) tea.Cmd {
	d.table.SetTable(table)
	return WorkspaceSet
}

func (d *Workspace) IsTableSet() bool {
	return d.table.IsTableSet()
}

func (d *Workspace) UpdateSize(termWidth, termHeight int) {
	panelWidth := termWidth/ui.PanelWidthRatio - ui.BorderThicknessDouble
	d.width, d.height = termWidth-panelWidth-ui.BorderThicknessDouble-1, termHeight-5
	// d.table.SetHeight(d.height)
}

func (d *Workspace) Style() lipgloss.Style {
	width := d.width
	if len(d.table.Table.Columns()) > 0 {
		width = d.table.Table.Width()
	}
	return workspaceStyle.Width(width).Height(d.height)
}

func (d *Workspace) View() string {
	// dataBorder := lipgloss.Border{
	// 	Top:         "─",
	// 	Bottom:      "─",
	// 	Right:       "│",
	// 	Left:        "│",
	// 	TopLeft:     "┬",
	// 	TopRight:    "┐",
	// 	BottomLeft:  "┴",
	// 	BottomRight: "┤",
	// }

	// width := d.width
	// if len(d.table.Table.Columns()) > 0 {
	// 	width = d.table.Table.Width()
	// }

	// workspaceStyle := lipgloss.NewStyle().
	// 	Width(width).
	// 	Height(d.height).
	// 	Border(dataBorder, true, true, false, true).
	// 	Align(lipgloss.Left, lipgloss.Top)

	if d.isLoading {
		lipgloss.Place(
			d.width,
			d.height,
			lipgloss.Center,
			lipgloss.Center,
			d.spinner.View(),
		)

		// return workspaceStyle.
		// 	AlignHorizontal(lipgloss.Center).
		// 	AlignVertical(lipgloss.Center).
		// 	Render(d.spinner.View())
	}

	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Left,
		lipgloss.Top,
		d.table.View(),
	)
	// tableView := d.table.View()

	// return workspaceStyle.Render(tableView)
}
