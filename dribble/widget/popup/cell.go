package popup

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/dribble/util"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
)

const scrollSpeed = 3

type TableCell struct {
	viewport.Model
	// MaxWidth, MaxHeight         int
	Value     string
	CancelCmd tea.Cmd
}

func newCellData(value string) *TableCell {
	value = util.PrettifyJson(value)
	// vp := viewport.New(lipgloss.Width(value), lipgloss.Height(value))
	vp := viewport.New(0, 0)
	vp.SetContent(value)

	return &TableCell{
		Model: vp,
		// MaxWidth:  maxWidth,
		// MaxHeight: maxHeight,
		Value:     value,
		CancelCmd: func() tea.Msg { return widget.PopupCancelMsg{} },
	}
}

// Exec implements PopupModel.
func (t *TableCell) Exec() tea.Cmd {
	return nil
}

// Cancel implements PopupModel.
func (t *TableCell) Cancel() tea.Cmd {
	return nil
}

// Init implements tea.Model.
func (t *TableCell) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *TableCell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Back, config.Keys.Quit, config.Keys.Select):
			return t, t.CancelCmd
		case key.Matches(msg, config.Keys.Up):
			t.SetYOffset(t.YOffset - scrollSpeed)
		case key.Matches(msg, config.Keys.Down):
			t.SetYOffset(t.YOffset + scrollSpeed)
		case key.Matches(msg, config.Keys.Left): // Maybe someday, since the XOffset is confusingly xOffset...
			// t.XOffset = t.XOffset + scrollSpeed
		case key.Matches(msg, config.Keys.Right): // Maybe someday
			// t.XOffset = t.XOffset + scrollSpeed
		}
	}
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)

	return t, cmd
}

// GetContentSize implements PopupModel.
func (t *TableCell) GetContentSize() (int, int) {
	return t.GetContentWidth(), t.GetContentHeight()
}

// GetContentWidth implements PopupModel.
func (t *TableCell) GetContentWidth() int {
	return lipgloss.Width(t.Value)
}

// GetContentHeight implements PopupModel.
func (t *TableCell) GetContentHeight() int {
	return lipgloss.Height(t.Value)
}

// SetMaxSize implements PopupModel.
func (t *TableCell) SetMaxSize(width int, height int) {
	t.Model.Width = min(width-ui.PopupStyle.GetHorizontalFrameSize(), t.GetContentWidth())
	t.Model.Height = min(height-ui.PopupStyle.GetVerticalFrameSize(), t.GetContentHeight())

	t.Model.SetContent(t.Value)
}

// View implements tea.Model.
func (t *TableCell) View() string {
	return t.Model.View()
}
