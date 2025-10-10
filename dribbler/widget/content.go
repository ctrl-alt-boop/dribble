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

		Style         lipgloss.Style
		LayoutManager layout.Manager

		Children     []tea.Model
		FocusedChild int
	}
)

func NewContentArea(id int, name string, children ...tea.Model) ContentArea {
	return ContentArea{
		ID:            id,
		name:          name,
		Children:      children,
		LayoutManager: &layout.SimpleLayout{},
	}
}

func New(name string, manager layout.Manager, children ...tea.Model) ContentArea {
	return ContentArea{
		name:          name,
		LayoutManager: manager,
		Children:      children,
		FocusedChild:  -1,
	}
}

func (a *ContentArea) Name() string {
	return a.name
}

func (a *ContentArea) AddChild(child ContentArea) {
	a.Children = append(a.Children, child)
}

func (a *ContentArea) SetLayoutManager(manager layout.Manager) {
	a.LayoutManager = manager
}

func (a *ContentArea) SetStyle(style lipgloss.Style) {
	a.Style = style
}

func (a ContentArea) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, child := range a.Children {
		cmds = append(cmds, child.Init())
	}
	return tea.Batch(cmds...)
}

func (a *ContentArea) UnfocusCmd() tea.Msg {
	return FocusBackMsg{}
}

func (a ContentArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	updated := a

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		updated.LayoutManager.SetSize(msg.Width, msg.Height)

		updatedChildren := a.LayoutManager.Layout(a.Children)
		updated.Children = updatedChildren
		return updated, nil

	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.CycleView) {
			updated.FocusedChild = (a.FocusedChild + 1) % len(a.Children)
			return updated, nil
		}
		if len(a.Children) == 0 {
			return updated, nil
		}
		focusedChild := a.Children[a.FocusedChild]
		focusedChild, cmd := focusedChild.Update(msg)
		updated.Children[a.FocusedChild] = focusedChild
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case FocusMsg:
		updated.FocusedChild = msg.Index[0]
		if updated.FocusedChild != -1 && updated.FocusedChild >= len(a.Children) {
			if len(msg.Index) > 1 {
				msg.Index = msg.Index[1:]
				focusedChild, cmd := a.Children[a.FocusedChild].Update(msg)
				updated.Children[a.FocusedChild] = focusedChild
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
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
	return a.Style.Render(a.LayoutManager.View(a.Children))
}
