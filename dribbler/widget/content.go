package widget

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui/layout"
)

var _ tea.Model = (*ContentArea)(nil)

type (
	ContentArea struct {
		ID   int
		name string

		Style  lipgloss.Style
		Layout layout.Manager

		Children     []tea.Model
		FocusedChild int
	}

	UnfocusChildMsg struct {
		ID int
	}
)

func NewContentArea(id int, name string, children ...tea.Model) ContentArea {
	return ContentArea{
		ID:       id,
		name:     name,
		Children: children,
		Layout:   &layout.SimpleLayout{},
	}
}

func (a *ContentArea) Name() string {
	return a.name
}

func (a *ContentArea) AddChild(child ContentArea) {
	a.Children = append(a.Children, child)
}

func (a *ContentArea) SetLayoutManager(manager layout.Manager) {
	a.Layout = manager
}

func (a *ContentArea) SetStyle(style lipgloss.Style) {
	a.Style = style
}

func (a ContentArea) Init() tea.Cmd {
	return nil
}

func (a *ContentArea) UnfocusCmd() tea.Msg {
	return UnfocusChildMsg{
		ID: a.FocusedChild,
	}
}

func (a ContentArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	updated := a

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		updated.Layout.SetSize(msg.Width, msg.Height)

		updatedChildren := a.Layout.Layout(a.Children)
		updated.Children = updatedChildren
		return updated, nil

	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Back) {
			return updated, a.UnfocusCmd
		} else if key.Matches(msg, config.Keys.CycleView) {
			updated.FocusedChild = (a.FocusedChild + 1) % len(a.Children)
			return updated, nil
		}
		focusedChild := a.Children[a.FocusedChild]
		focusedChild, cmd := focusedChild.Update(msg)
		updated.Children[a.FocusedChild] = focusedChild
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	default:
		for i, child := range a.Children {
			child, cmd := child.Update(msg)
			updated.Children[i] = child
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return updated, tea.Batch(cmds...)
}

func (a ContentArea) View() string {
	return a.Style.Render(a.Layout.View(a.Children))
}
