package components

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribbler/keys"
)

type (
	// Popup orchestrates popups
	Popup struct {
		Base
	}

	// PopupCloseMsg is used to command the popup to be closed
	PopupCloseMsg struct{}
)

// Update implements tea.Model.
func (p *Popup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Map.Back):
			return p, func() tea.Msg {
				return PopupCloseMsg{}
			}
		}
	}

	return p, nil
}

// // Update implements tea.Model.
// func (t *TableCell) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, keys.Keys.Back, keys.Keys.Quit, keys.Keys.Select):
// 			return t, t.CancelCmd
// 		case key.Matches(msg, keys.Keys.Up):
// 			t.SetYOffset(t.YOffset - scrollSpeed)
// 		case key.Matches(msg, keys.Keys.Down):
// 			t.SetYOffset(t.YOffset + scrollSpeed)
// 		case key.Matches(msg, keys.Keys.Left): // Maybe someday, since the XOffset is confusingly xOffset...
// 			// t.XOffset = t.XOffset + scrollSpeed
// 		case key.Matches(msg, keys.Keys.Right): // Maybe someday
// 			// t.XOffset = t.XOffset + scrollSpeed
// 		}
// 	}
// 	var cmd tea.Cmd
// 	t.Model, cmd = t.Model.Update(msg)

// 	return t, cmd
// }

// // Init implements PopupModel.
// func (c *Connect) Init() tea.Cmd {
// 	formTitle := fmt.Sprintf("Connect with %s driver", c.driverName)
// 	c.ipInput = huh.NewInput().
// 		Key("ip").
// 		Value(&c.ip).
// 		Title("IP:")
// 	c.portInput = huh.NewInput().
// 		Key("port").
// 		Value(&c.port).
// 		Title("Port:")

// 	c.form = huh.NewForm(
// 		huh.NewGroup(
// 			huh.NewNote().Title(formTitle),

// 			huh.NewInput().Title("Username:").
// 				Value(&c.username).
// 				Validate(notEmpty),
// 			huh.NewInput().Title("Password:").
// 				Value(&c.password).
// 				EchoMode(huh.EchoModePassword),
// 			huh.NewConfirm().Title("Default server settings?").
// 				Key("defaultServer").
// 				Value(&c.defaultServer).
// 				Affirmative("Y").
// 				Negative("N"),
// 		),
// 		huh.NewGroup(
// 			c.ipInput,
// 			c.portInput).
// 			WithHideFunc(func() bool {
// 				return c.defaultServer
// 			}),
// 		// huh.NewGroup(
// 		// 	huh.NewConfirm().Title("Login?").
// 		// 		Key("confirm").
// 		// 		Affirmative("Y").
// 		// 		Negative("N").
// 		// 		Value(&c.confirm),
// 		// ),
// 	).
// 		WithLayout(huh.LayoutStack).
// 		WithShowHelp(false).
// 		WithShowErrors(false)
// 	// WithWidth(width).WithHeight(height)

// 	c.stableWidth = lipgloss.Width(c.form.View())
// 	c.stableHeight = lipgloss.Height(c.form.View())
// 	return c.form.Init()
// }

// // Update implements PopupModel.
// func (c *Connect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	form, cmd := c.form.Update(msg)
// 	if f, ok := form.(*huh.Form); ok {
// 		c.form = f
// 	}

// 	switch c.form.State {
// 	case huh.StateCompleted:
// 		return c, c.ConfirmCmd
// 	case huh.StateAborted:
// 		return c, c.CancelCmd
// 	}

// 	keyMsg, ok := msg.(tea.KeyMsg)
// 	if ok && key.Matches(keyMsg, keys.Keys.Back) {
// 		return c, c.CancelCmd
// 	}

// 	return c, cmd
// }
