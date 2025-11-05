package explorer

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const sidePanelID string = string(ExplorerPageID) + ".sidepanel"

const numberSelectionDuration = 500 * time.Millisecond

type selectMsg int

type sidePanelMode int

const (
	modeList sidePanelMode = iota
	modeTree
)

type sidePanelItem struct { // Maybe I'll change this to some interface with OnSelectCmd() tea.Cmd
	onSelectCmd tea.Cmd
}

func (s sidePanelItem) FilterValue() string { return "" }
func (s sidePanelItem) Title() string       { return "" }
func (s sidePanelItem) Description() string { return "" }

type sidePanel struct {
	mode sidePanelMode
	list list.Model

	innerWidth, innerHeight int

	inNumberSelection bool
	numberSelection   string
}

func newSidePanel() *sidePanel {
	return &sidePanel{
		innerWidth:  0,
		innerHeight: 0,

		mode: modeList,

		inNumberSelection: false,
		numberSelection:   "",
	}
}

func (s *sidePanel) GetSelected() any {
	return nil
}

func (s *sidePanel) SetInnerSize(width, height int) {
	s.innerWidth, s.innerHeight = width, height
}

func (s *sidePanel) Init() tea.Cmd {
	itemDelegate := list.NewDefaultDelegate()
	itemDelegate.ShowDescription = false
	itemDelegate.SetSpacing(1)
	itemDelegate.SetHeight(1)

	itemDelegate.Styles.SelectedTitle = lipgloss.NewStyle().Bold(true)
	itemDelegate.Styles.SelectedDesc = lipgloss.NewStyle().Bold(true)
	itemDelegate.Styles.NormalTitle = lipgloss.NewStyle().Faint(true)
	itemDelegate.Styles.NormalDesc = lipgloss.NewStyle().Faint(true)

	s.list = list.New([]list.Item{}, itemDelegate, s.innerWidth, s.innerHeight)

	return nil
}

func (s *sidePanel) Update(msg tea.Msg) (*sidePanel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Up):
			s.list.CursorUp()

		case key.Matches(msg, keyMap.Down):
			s.list.CursorDown()

		case key.Matches(msg, keyMap.Select):
			if s.inNumberSelection {
				selectIndex, err := strconv.Atoi(s.numberSelection)
				if err != nil {
					return s, nil
				}
				s.list.Select(selectIndex)
				s.inNumberSelection = false
				s.numberSelection = ""
			}
			selected := s.list.SelectedItem()
			return s, onSelect(selected)

		case key.Matches(msg, keyMap.Number):
			s.inNumberSelection = true
			s.numberSelection += msg.String()
			return s, s.debounceSelection()

		case key.Matches(msg, keyMap.Back):
			if s.inNumberSelection && len(s.numberSelection) > 0 {
				s.numberSelection = s.numberSelection[:len(s.numberSelection)-1]
				if len(s.numberSelection) == 0 {
					s.inNumberSelection = false
					return s, nil
				}
				return s, s.debounceSelection()
			}

		case key.Matches(msg, keyMap.Esc):
			if s.inNumberSelection {
				s.inNumberSelection = false
				s.numberSelection = ""
				return s, nil
			}
		}
	}
	return s, nil
}

const (
	numSelectionLines = 1
)

func (s *sidePanel) Render() *lipgloss.Layer {
	selectionBox := s.createNumSelectionBox()

	s.list.SetSize(s.innerWidth, s.innerHeight-lipgloss.Height(selectionBox))

	view := lipgloss.JoinVertical(lipgloss.Left, s.list.View())
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		Width(s.innerWidth).Height(s.innerHeight)
	view = box.Render(view)

	return lipgloss.NewLayer(view)
}

func (s *sidePanel) createNumSelectionBox() string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		Width(s.innerWidth).Height(numSelectionLines+lipgloss.NormalBorder().GetTopSize()).
		PaddingLeft(1).
		Render("sel: ", s.numberSelection)
}

func (s *sidePanel) debounceSelection() tea.Cmd {
	return tea.Tick(numberSelectionDuration, func(_ time.Time) tea.Msg {
		if !s.inNumberSelection {
			return nil
		}
		defer func() {
			s.numberSelection = ""
			s.inNumberSelection = false
		}()
		selectIndex, err := strconv.Atoi(s.numberSelection)
		if err != nil {
			return nil
		}
		return selectMsg(selectIndex)
	})
}

func onSelect(item list.Item) tea.Cmd {
	return item.(*sidePanelItem).onSelectCmd
}
