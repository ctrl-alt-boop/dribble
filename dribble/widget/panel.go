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
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
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

type Panel struct {
	list *ui.List

	Width, Height           int
	InnerWidth, InnerHeight int

	showDetails bool

	goolDb *gooldb.GoolDb

	mode      PanelMode
	isLoading bool
	spinner   spinner.Model

	selectIndexHistory []int
}

func NewPanel(gool *gooldb.GoolDb) *Panel {
	return &Panel{
		list:               ui.NewList(),
		goolDb:             gool,
		mode:               ServerList,
		spinner:            spinner.New(spinner.WithSpinner(ui.MovingBlock)),
		selectIndexHistory: make([]int, 0),
		showDetails:        config.Cfg.Ui.ShowDetails,
	}
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
		selection, ok := p.list.SelectedItem().(*ui.ConnectionItem)
		if ok {
			cmd = func() tea.Msg {
				return SelectServerMsg(string(selection.Name))
			}
		}
	case DatabaseList:
		selection, ok := p.list.SelectedItem().(*ui.ConnectionItem)
		if ok {
			p.isLoading = true
			cmd = func() tea.Msg {
				return SelectDatabaseMsg(string(selection.Name))
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
	var connectionItems []*ui.ConnectionItem
	connectionItems = append(connectionItems, ui.GetSavedConfigsSorted()...)

	for name, settings := range config.GetDriverDefaults() {
		connectionItems = append(connectionItems, &ui.ConnectionItem{
			Name: name,
			Settings: connection.NewSettings(connection.AsType(connection.Driver),
				connection.WithDriver(settings.DriverName),
				connection.WithHost(settings.Ip, settings.Port),
				connection.WithUser(settings.Username),
				connection.WithPassword(settings.Password)),
		})
	}
	p.list.SetConnectionItems(connectionItems)
	return nil // should maybe do this in AppModel with a Cmd
}

func (p *Panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Up):
			p.list.CursorUp()
		case key.Matches(msg, config.Keys.Down):
			p.list.CursorDown()
		case key.Matches(msg, config.Keys.Select):
			return p, p.OnSelect()
		case key.Matches(msg, config.Keys.Back):
			return p, nil
		case key.Matches(msg, config.Keys.Details):
			p.showDetails = !p.showDetails
			return p, nil
		}

	case io.GoolDbEventMsg:
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
				items := ui.SettingsToConnectionItems(args.Databases)
				p.list.SetConnectionItems(items)
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

func (p *Panel) UpdateSize(width, height int) {
	// p.width, p.height = width/ui.PanelWidthRatio-ui.BorderThicknessDouble, height-5
	p.Width, p.Height = width, height
	p.InnerWidth = p.Width - ui.PanelStyle.GetHorizontalFrameSize()
	p.InnerHeight = p.Height - ui.PanelStyle.GetVerticalFrameSize()
}

func (p *Panel) View() string {
	if p.isLoading {
		return lipgloss.Place(
			p.InnerWidth,
			p.InnerHeight,
			lipgloss.Center,
			lipgloss.Position(0.2),
			p.spinner.View(),
		)
	}

	details, ok := p.list.SelectedItem().(ui.ConnectionItem)
	var detailsView string
	if p.showDetails && ok {
		detailsView = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).BorderTop(true).
			Padding(1, 0).
			MaxHeight(9).Width(p.InnerWidth).
			Render(details.Inspect())
	}

	p.list.SetSize(p.InnerWidth, p.InnerHeight-lipgloss.Height(detailsView))

	joined := lipgloss.JoinVertical(lipgloss.Left, p.list.View(), detailsView)

	return ui.PanelStyle.Width(p.InnerWidth).Height(p.InnerHeight).Render(joined)
}
