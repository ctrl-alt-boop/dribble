package explorer

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var empty = &emptyTab{}

type emptyTab struct{}

// Title implements Tab.
func (e *emptyTab) Name() string { return "" }

// Update implements Tab.
func (e emptyTab) Update(_ tea.Msg) tea.Cmd { return nil }

// Render implements Tab.
func (e emptyTab) Render() string {
	return ""
}

type Tab interface {
	Name() string

	Update(msg tea.Msg) tea.Cmd

	Render() string
}

type Tabs struct {
	names   []string
	tabs    []Tab
	current int

	normalStyle  lipgloss.Style
	currentStyle lipgloss.Style
}

func NewTabs() *Tabs {
	return &Tabs{
		names:   []string{},
		tabs:    []Tab{},
		current: 0,

		normalStyle:  unstyled.Border(lipgloss.NormalBorder()).Padding(0, 1),
		currentStyle: unstyled.Border(lipgloss.ThickBorder()).Padding(0, 1),
	}
}

func (t *Tabs) Add(tab Tab) {
	t.tabs = append(t.tabs, tab)
	t.names = append(t.names, tab.Name())
	t.current = len(t.tabs) - 1
}

func (t *Tabs) Set(tabs ...Tab) {
	t.tabs = tabs
	t.names = []string{}
	for _, tab := range tabs {
		t.names = append(t.names, tab.Name())
	}
}

func (t *Tabs) Close() {
	if len(t.tabs) == 0 {
		return
	}
	t.tabs = append(t.tabs[:t.current], t.tabs[t.current+1:]...)
	t.names = append(t.names[:t.current], t.names[t.current+1:]...)
	if t.current >= len(t.tabs) {
		t.current = len(t.tabs) - 1
	} else {
		t.current--
	}
}

func (t Tabs) Current() Tab {
	if t.tabs == nil {
		return empty
	}
	return t.tabs[t.current]
}

func (t *Tabs) Move(v int) {
	if t.tabs == nil {
		return
	}
	t.current = (t.current + v + len(t.tabs)) % len(t.tabs)
}

func (t Tabs) Len() int {
	return len(t.tabs)
}

func (t *Tabs) Clear() {
	t.tabs = []Tab{}
	t.names = []string{}
	t.current = 0
}

func (t Tabs) IsEmpty() bool {
	return len(t.tabs) == 0
}

func (t *Tabs) UpdateCurrent(msg tea.Msg) tea.Cmd {
	return t.Current().Update(msg)
}

func (t *Tabs) UpdateAll(msg tea.Msg) tea.Cmd {
	if t.tabs == nil {
		return nil
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	for _, tab := range t.tabs {
		cmd = tab.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (t Tabs) Render() (string, string) {
	if t.tabs == nil {
		return "", ""
	}

	return t.Current().Render(), t.RenderTabs()
}

func (t Tabs) RenderTabs() string {
	tabNames := []string{}
	if len(t.names) == 0 {
		return ""
	}
	widestTab := 0
	for _, name := range t.names {
		if len(name) > widestTab {
			widestTab = len(name)
		}
	}
	tabWidth := widestTab + t.currentStyle.GetHorizontalFrameSize()

	for i, name := range t.names {
		if i == t.current {
			tabNames = append(tabNames, getTabStyle(tabWidth, true).Render(name))
		} else {
			tabNames = append(tabNames, getTabStyle(tabWidth, false).Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabNames...)
}

type StringTab struct {
	name    string
	content string
}

func (s StringTab) Name() string {
	return s.name
}

func (s StringTab) Render() string {
	return s.content
}

func (s StringTab) Update(msg tea.Msg) tea.Cmd {
	return nil
}
