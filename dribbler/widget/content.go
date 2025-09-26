package widget

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/widget/layout"
)

var _ tea.Model = (*ContentArea)(nil)

type (
	ContentArea struct {
		ID   int
		name string

		viewport viewport.Model

		Style  lipgloss.Style
		Layout layout.Manager

		Width, Height int
		X, Y          int

		Children     []tea.Model
		FocusedChild int
	}

	UnfocusChildMsg struct {
		ID int
	}
)

func NewContentArea(id int, name string) *ContentArea {
	return &ContentArea{
		ID:       id,
		name:     name,
		viewport: viewport.New(0, 0),
		Layout:   &layout.SimpleLayout{},
	}
}

func (a *ContentArea) Name() string {
	return a.name
}

func (a *ContentArea) AddChild(child *ContentArea) {
	a.Children = append(a.Children, child)
}

func (a *ContentArea) SetSyle(style lipgloss.Style) {
	a.Style = style
}

func (a *ContentArea) UpdateSize(width int, height int) {
	a.Width, a.Height = width, height
}

func (a *ContentArea) InnerSize() (int, int) {
	return a.Width - a.Style.GetHorizontalFrameSize(), a.Height - a.Style.GetVerticalFrameSize()
}

func (a *ContentArea) Init() tea.Cmd {
	return nil
}

func (a *ContentArea) UnfocusCmd() tea.Msg {
	return UnfocusChildMsg{
		ID: a.FocusedChild,
	}
}

func (a *ContentArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.Width, a.Height = msg.Width, msg.Height
		cmd := a.Layout.Layout(a.Width, a.Height, a.Children)

		return a, cmd
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Back) {
			return a, a.UnfocusCmd
		} else if key.Matches(msg, config.Keys.CycleView) {
			a.FocusedChild = (a.FocusedChild + 1) % len(a.Children)
			return a, nil
		}
		focusedChild := a.Children[a.FocusedChild]
		focusedChild, cmd := focusedChild.Update(msg)
		a.Children[a.FocusedChild] = focusedChild
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	default:
		for i, child := range a.Children {
			child, cmd := child.Update(msg)
			a.Children[i] = child
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return a, tea.Batch(cmds...)
}

func (a *ContentArea) View() string {
	a.viewport.Width = a.Width
	a.viewport.Height = a.Height

	a.viewport.SetContent(a.Layout.View(a.Children))

	return a.Style.Width(a.Width).Height(a.Height).Render(a.viewport.View())
}
