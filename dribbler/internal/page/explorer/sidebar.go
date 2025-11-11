package explorer

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribbler/internal/config"
	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"

	"github.com/ctrl-alt-boop/dribbler/logging"
)

const sidebarID string = string(ExplorerPageID) + ".sidebar"

const numberSelectionDuration = 1000 * time.Millisecond

type debounceSelectMsg int

type sidebarMode int

const (
	modeList sidebarMode = iota
	modeTree
)

type SidebarItem interface {
	list.DefaultItem
	OnSelect() tea.Msg
}

type sidebar struct {
	mode        sidebarMode
	list        list.Model
	searchInput textinput.Model

	width, height           int
	innerWidth, innerHeight int

	shouldDebounceSelect bool
	inNumberSelection    bool
	numberSelection      string

	showDetails bool

	style lipgloss.Style

	keybinds *SidebarKeys
}

func newSidebar() *sidebar {
	return &sidebar{
		innerWidth:  0,
		innerHeight: 0,

		mode: modeList,

		shouldDebounceSelect: false,
		inNumberSelection:    false,
		numberSelection:      "",

		showDetails: false,

		keybinds: DefaultSidebarKeyBindings(),
	}
}

func (s *sidebar) GetSelected() SidebarItem {
	return s.list.SelectedItem().(SidebarItem)
}

func (s *sidebar) SetMode(mode sidebarMode) {
	s.mode = mode
}

func (s *sidebar) SetShouldDebounceSelect(shouldDebounceSelect bool) {
	s.shouldDebounceSelect = shouldDebounceSelect
}

func (s *sidebar) SetItems(items []SidebarItem) {
	itms := make([]list.Item, 0, len(items))
	for _, item := range items {
		items = append(items, item)
	}
	s.list.SetItems(itms)
}

func (s *sidebar) AddItem(item SidebarItem) {
	s.list.InsertItem(len(s.list.Items()), item) // len items should append it per the InsertItem dok (if index >= len = append)
}

func (s *sidebar) SetSize(width, height int) {
	s.width, s.height = width, height
	s.innerWidth, s.innerHeight = width-s.style.GetHorizontalFrameSize(), height-s.style.GetVerticalFrameSize()
}

func (s *sidebar) Init() tea.Cmd {
	logging.GlobalLogger().Infof("sidebar.Init")

	s.initList()

	saved := s.createSavedConfigList()
	supported := s.createSupportedSourcesList()

	sidebarItems := slices.Concat(saved, supported)
	s.list.SetItems(sidebarItems)

	return nil
}

func (s *sidebar) initList() {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)
	delegate.SetHeight(1)
	// delegate.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
	// 	if msg, ok := msg.(tea.KeyMsg); ok {
	// 		if key.Matches(msg, s.keybinds.Accept) {
	// 			logging.GlobalLogger().Infof("Accept: %d", m.GlobalIndex())
	// 			return onSelect(m.Items()[m.GlobalIndex()])
	// 		}
	// 	}
	// 	return nil
	// }

	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Bold(true)
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Faint(true)
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Faint(true)

	s.list = list.New([]list.Item{}, delegate, s.innerWidth, s.innerHeight)

	s.list.SetShowTitle(false)
	s.list.SetShowStatusBar(false)
	s.list.SetShowHelp(false)
	s.list.SetShowPagination(false)
	s.list.SetShowFilter(false)
	// s.list.SetFilteringEnabled(false)

	s.list.KeyMap = list.DefaultKeyMap()
	s.list.DisableQuitKeybindings()

	s.list.Filter = list.DefaultFilter

	s.searchInput = textinput.New()
	s.searchInput.Prompt = "/"
	s.searchInput.CharLimit = 16
	s.list.FilterInput = s.searchInput
}

func (s *sidebar) Update(msg tea.Msg) (*sidebar, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		logging.GlobalLogger().Infof("sidebar.KeyPressMsg: %s", msg)
		switch {
		case key.Matches(msg, s.keybinds.Accept):
			if s.inNumberSelection {
				selectIndex, err := strconv.Atoi(s.numberSelection)
				s.inNumberSelection = false
				s.numberSelection = ""
				if err != nil {
					logging.GlobalLogger().Warnf("Atoi err: %v", err)
					return s, nil
				}
				selectIndex--
				if selectIndex < 0 || selectIndex >= len(s.list.Items()) {
					logging.GlobalLogger().Infof("selectIndex %d out of bounds: length items: %d", selectIndex, len(s.list.Items()))
					return s, nil
				}
				s.list.Select(selectIndex)
			}
			selected := s.list.SelectedItem()
			cmd = s.onSelect(selected)

		case key.Matches(msg, s.keybinds.Number):
			s.inNumberSelection = true
			s.numberSelection += msg.String()
			cmd = s.debounceSelection()

		case key.Matches(msg, s.keybinds.Cancel):
			if s.inNumberSelection && len(s.numberSelection) > 0 {
				s.numberSelection = s.numberSelection[:len(s.numberSelection)-1]
				if len(s.numberSelection) == 0 {
					s.inNumberSelection = false
					return s, nil
				}
				cmd = s.debounceSelection()
			}

		case key.Matches(msg, s.keybinds.Details):
			s.showDetails = !s.showDetails

		case key.Matches(msg, s.keybinds.Search):
			s.list.SetShowFilter(true)

		case key.Matches(msg, s.keybinds.Esc):
			if s.inNumberSelection {
				s.inNumberSelection = false
				s.numberSelection = ""
				return s, nil
			}
			if s.list.FilterState() == list.Filtering {
				s.list.SetShowFilter(false)
			}
		}

	case debounceSelectMsg:
		if msg < 0 || int(msg) >= len(s.list.Items()) {
			return s, nil
		}
		s.list.Select(int(msg))
		return s, s.onSelect(s.list.SelectedItem())

	case dribbleapi.DribbleResponseMsg:
		logging.GlobalLogger().Infof("sidebar.DribbleResponseMsg: %+v", msg.Response)
		// switch msg := msg.Response.(type) {
		// case dribbleapi.TargetOpenedMsg:
		// }

	}
	var listCmd tea.Cmd
	s.list, listCmd = s.list.Update(msg)

	return s, tea.Batch(cmd, listCmd)
}

const (
	numSelectionLines = 1
)

func (s *sidebar) SetStyle(style lipgloss.Style) {
	s.style = style
}

func (s *sidebar) Render() *lipgloss.Layer {
	var selectionBox string
	if s.inNumberSelection {
		selectionBox = s.createNumSelectionBox()
	}

	s.list.SetSize(s.innerWidth, s.innerHeight-lipgloss.Height(selectionBox))

	view := lipgloss.JoinVertical(lipgloss.Left, s.list.View(), selectionBox)
	style := s.style.
		Width(s.width).Height(s.height)
	view = style.Render(view)

	return lipgloss.NewLayer(view)
}

func (s *sidebar) createNumSelectionBox() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		Width(s.innerWidth).Height(numSelectionLines+lipgloss.NormalBorder().GetTopSize()).
		PaddingLeft(1).
		Render("sel: ", s.numberSelection)
}

func (s *sidebar) debounceSelection() tea.Cmd {
	if !s.shouldDebounceSelect {
		return nil
	}
	return tea.Tick(numberSelectionDuration, func(_ time.Time) tea.Msg {
		if !s.inNumberSelection {
			return nil
		}
		s.inNumberSelection = false

		selectIndex, err := strconv.Atoi(s.numberSelection)
		s.numberSelection = ""
		if err != nil {
			logging.GlobalLogger().Warnf("Atoi err: %v", err)
			return nil
		}
		selectIndex--
		if selectIndex < 0 || selectIndex >= len(s.list.Items()) {
			logging.GlobalLogger().Infof("selectIndex %d out of bounds: length items: %d", selectIndex, len(s.list.Items()))
			return nil
		}
		return debounceSelectMsg(selectIndex)
	})
}

func (s *sidebar) createSavedConfigList() []list.Item {
	items := make([]list.Item, 0, len(config.SavedConfigs))
	for confName, dsn := range config.SavedConfigs {
		items = append(items, &DataSourceItem{
			Name: confName,
			Properties: &DataSourceItemProps{
				SourceType: string(dsn.SourceType()),
				Type:       "saved",
				Info:       dsn.Info(),
				Extras: map[string]string{
					"DSN": dsn.DSN(),
				},
				createdFrom: "datasource.Namer",
			},
			object: dsn,
		})
	}
	return items
}

func (s *sidebar) createSupportedSourcesList() []list.Item {
	supportedAdapters := datasource.Adapters()
	items := make([]list.Item, 0, len(supportedAdapters))
	for _, adapter := range supportedAdapters {
		items = append(items, &DataSourceItem{
			Name: adapter.Name,
			Properties: &DataSourceItemProps{
				SourceType:  string(adapter.Type),
				Type:        string(adapter.Metadata.StorageType),
				Info:        strings.Join(adapter.Capabilities.AsStrings(), ", "),
				Extras:      adapter.Properties,
				createdFrom: "datasource.Adapter",
			},
			object: adapter,
		})
	}
	return items
}

func (s *sidebar) onSelect(item list.Item) tea.Cmd {
	switch item := item.(type) {
	case *DataSourceItem:
		if s.showDetails {
			logging.GlobalLogger().Infof("%#v", item)
		} else {
			logging.GlobalLogger().Infof("%s", item)
		}
		return item.OnSelect
	}
	logging.GlobalLogger().Warnf("selected item not handled: %T: %+v", item, item)
	return nil
}

var _ SidebarItem = (*DataSourceItem)(nil)

type DataSourceItem struct {
	Name       string
	Properties *DataSourceItemProps

	object any
}

func (i DataSourceItem) String() string {
	if i.Properties == nil {
		return i.Name
	}
	return fmt.Sprintf("%s: %s", i.Name, i.Properties)
}

// GoString implements fmt.GoStringer.
func (i *DataSourceItem) GoString() string {
	extras := []string{}
	for k, v := range i.Properties.Extras {
		extras = append(extras, fmt.Sprintf("%s: %s", k, v))
	}
	extraString := fmt.Sprintf("{%s}", strings.Join(extras, ", "))
	internals := fmt.Sprintf("%s, %s, %s, %s (%s)", i.Properties.SourceType, i.Properties.Info, extraString, i.Properties.createdFrom, i.Properties.Type)
	return fmt.Sprintf("%s: { %s }", i.Name, internals)
}

// Description implements SidebarItem.
func (i *DataSourceItem) Description() string {
	return fmt.Sprintf("%s", i.Properties)
}

// OnSelect implements SidebarItem.
func (i *DataSourceItem) OnSelect() tea.Msg {
	switch object := i.object.(type) {
	case datasource.Adapter:
		return dribbleapi.AdapterOpen(i.Name, object)
	case datasource.Namer:
		return dribbleapi.DSNOpen(i.Name, object)
	}
	logging.GlobalLogger().Warnf("OnSelect not handled: %T: %+v", i.object, i.object)
	return nil
}

// Title implements SidebarItem.
func (i *DataSourceItem) Title() string {
	return i.Name
}

// FilterValue implements list.Item.
func (i *DataSourceItem) FilterValue() string {
	v := i.Name
	if i.Properties != nil {
		v += " " + i.Properties.SourceType
		v += " " + i.Properties.Info
		v += " " + i.Properties.createdFrom
		v += " " + i.Properties.Type
		for _, v := range i.Properties.Extras {
			v += " " + v
		}

	}

	return v
}

type DataSourceItemProps struct {
	Type        string
	SourceType  string
	Info        string
	Extras      map[string]string
	createdFrom string
}

func (p DataSourceItemProps) String() string {
	return fmt.Sprintf("%s: %s", p.SourceType, p.Info)
}

func (s *sidebar) CreateListKeyMap() list.KeyMap {
	return list.KeyMap{
		CursorUp:             s.keybinds.Up,
		CursorDown:           s.keybinds.Down,
		Filter:               s.keybinds.Search,
		ClearFilter:          s.keybinds.Esc,
		CancelWhileFiltering: s.keybinds.Esc,
		AcceptWhileFiltering: s.keybinds.Accept,
	}
}
