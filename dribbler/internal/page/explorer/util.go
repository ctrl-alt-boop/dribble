package explorer

import (
	"container/ring"

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
func (e emptyTab) Render() *lipgloss.Layer {
	return lipgloss.NewLayer("")
}

type Tab interface {
	Name() string

	Render() *lipgloss.Layer
	Update(msg tea.Msg) tea.Cmd
}

type Tabs struct {
	tabs *ring.Ring
}

func NewTabs() *Tabs {
	return &Tabs{
		tabs: nil,
	}
}

func (t *Tabs) Add(tab Tab) {
	if t.tabs == nil {
		t.tabs = ring.New(1)
		t.tabs.Value = tab

	} else {
		t.tabs = t.tabs.Link(&ring.Ring{Value: tab})
	}
}

func (t *Tabs) Set(tabs ...Tab) {
	t.tabs = ring.New(len(tabs))
	for _, tab := range tabs {
		t.tabs.Value = tab
		t.tabs = t.tabs.Next()
	}
}

func (t *Tabs) Close() {
	if t.tabs == nil {
		return
	}
	t.tabs = t.tabs.Prev()
	t.tabs = t.tabs.Unlink(1)
}

func (t *Tabs) Current() Tab {
	if t.tabs == nil {
		return empty
	}
	return t.tabs.Value.(Tab)
}

func (t *Tabs) Move(v int) {
	if t.tabs == nil {
		return
	}
	t.tabs = t.tabs.Move(v)
}

func (t *Tabs) Len() int {
	if t.tabs == nil {
		return 0
	}
	return t.tabs.Len()
}

func (t *Tabs) Clear() {
	t.tabs = nil
}

func (t *Tabs) IsEmpty() bool {
	return t.tabs == nil
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
	t.tabs.Do(
		func(a any) {
			cmd = a.(Tab).Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		},
	)
	return tea.Batch(cmds...)
}

func (t *Tabs) Render() (current *lipgloss.Layer, tabs *lipgloss.Layer) {
	if t.tabs == nil {
		return lipgloss.NewLayer(""), lipgloss.NewLayer("")
	}

	return t.tabs.Value.(Tab).Render(), lipgloss.NewLayer("")
}
