package widget

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/io"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type PanelSelectMsg struct {
	CurrentMode PanelMode
	Selected    string
}

type PanelMode string

const (
	ServerList   PanelMode = "serverList"
	DatabaseList PanelMode = "databaseList"
	TableList    PanelMode = "tableList"
)

type selection struct {
	index     int
	cachePath string
}

type Panel struct {
	list          *ui.List
	width, height int
	goolDb        *gooldb.GoolDb

	mode PanelMode

	isLoading bool

	spinner spinner.Model

	selectHistory []selection
}

func NewPanel(gool *gooldb.GoolDb) *Panel {

	return &Panel{
		list:          ui.NewList(),
		goolDb:        gool,
		mode:          ServerList,
		spinner:       spinner.New(spinner.WithSpinner(ui.MovingBlock)),
		selectHistory: make([]selection, 0),
	}
}

func (p *Panel) UpdateSize(termWidth, termHeight int) {
	p.width, p.height = termWidth/ui.PanelWidthRatio-ui.BorderThicknessDouble, termHeight-5
	p.list.SetSize(p.width, p.height)
}

func (p *Panel) SetMode(mode PanelMode) {
	p.mode = mode
}

func (p *Panel) GetMode() PanelMode {
	return p.mode
}

func (p *Panel) OnSelect() tea.Cmd {
	var cmd tea.Cmd
	switch p.mode {
	case ServerList:
		selection, ok := p.list.SelectedItem().(ui.ConnectionItem)
		if ok {
			logger.Infof("Selected: %+v", selection)
			cmd = func() tea.Msg {
				return SelectServerMsg(string(selection.Name))
			}
		}
	case DatabaseList:
		selection, ok := p.list.SelectedItem().(ui.ListItem)
		if ok {
			p.isLoading = true
			cmd = func() tea.Msg {
				return SelectDatabaseMsg(string(selection))
			}
		}

	case TableList:
		selection, ok := p.list.SelectedItem().(ui.ListItem)
		if ok {
			p.isLoading = true
			cmd = func() tea.Msg {
				return SelectTableMsg(string(selection))
			}
		}

	}

	return tea.Batch(cmd, p.spinner.Tick)
}

func (p *Panel) GetSelected() string {
	selection, ok := p.list.SelectedItem().(ui.ListItem)
	if !ok {
		return ""
	}
	return string(selection)
}

func (p *Panel) Select() tea.Msg {
	selection, ok := p.list.SelectedItem().(ui.ListItem)
	if !ok {
		return nil
	}
	return PanelSelectMsg{
		CurrentMode: p.mode,
		Selected:    string(selection),
	}
}

func (p *Panel) Init() tea.Cmd {
	var connectionItems []ui.ConnectionItem
	for name, settings := range config.SavedConfigs {
		connectionItems = append(connectionItems, ui.ConnectionItem{
			Name:     name,
			Settings: settings,
		})
	}
	for name, settings := range p.goolDb.GetDriverDefaults() {
		connectionItems = append(connectionItems, ui.ConnectionItem{
			Name:     name,
			Settings: settings,
		})
	}
	p.list.SetConnectionItems(connectionItems)
	return nil // should maybe do this in AppModel with a Cmd
}

func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Infof("Got message: %+v", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.UpdateSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Up):
			p.list.CursorUp()
		case key.Matches(msg, config.Keys.Down):
			p.list.CursorDown()
		case key.Matches(msg, config.Keys.Select):
			return p, p.OnSelect()
		}

	case io.GoolDbEventMsg:
		logger.Infof("Got GoolDbEventMsg: %+v", msg)
		p.isLoading = false
		p.spinner = spinner.New(spinner.WithSpinner(ui.MovingBlock))
		if msg.Err != nil {
			logger.Error(msg.Err)
			return p, nil
		}
		switch msg.Type {
		case gooldb.DatabaseListFetched:
			args, ok := msg.Args.(gooldb.DatabaseListFetchData)
			if ok {
				p.list.SetStringItems(args.Databases)
				p.SetMode(DatabaseList)
			}
		case gooldb.DBTableListFetched:
			args, ok := msg.Args.(gooldb.DBTableListFetchData)
			if ok {
				p.list.SetStringItems(args.Tables)
				p.SetMode(TableList)
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	}

	return p, nil
}

func (p *Panel) View() string {
	panelBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┬",
		BottomLeft:  "├",
		BottomRight: "┴",
	}

	panelStyle := lipgloss.NewStyle().
		Height(p.height).
		Width(p.width).
		Border(panelBorder, true, false, false, true).
		Align(lipgloss.Left, lipgloss.Top)

	if p.isLoading {
		return panelStyle.
			Height(p.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Top).
			Render(p.spinner.View())
	}

	return panelStyle.Render(p.list.View())
}
