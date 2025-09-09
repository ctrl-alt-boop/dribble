package popup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

type Connect struct {
	form      *huh.Form
	ipInput   *huh.Input
	portInput *huh.Input

	username string
	password string

	defaultServer bool
	ip            string
	port          string

	driverName string

	CancelCmd tea.Cmd

	stableWidth, stableHeight int
}

func notEmpty(s string) error {
	if s == "" {
		return errors.New("field cannot be empty")
	}
	return nil
}

func newConnect(s string) *Connect {
	connect := &Connect{
		defaultServer: true,
		driverName:    s,
		CancelCmd:     func() tea.Msg { return widget.PopupCancelMsg{} },
	}

	return connect
}

// Init implements PopupModel.
func (c *Connect) Init() tea.Cmd {
	formTitle := fmt.Sprintf("Connect with %s driver", c.driverName)
	c.ipInput = huh.NewInput().
		Key("ip").
		Value(&c.ip).
		Title("IP:")
	c.portInput = huh.NewInput().
		Key("port").
		Value(&c.port).
		Title("Port:")

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(formTitle),

			huh.NewInput().Title("Username:").
				Value(&c.username).
				Validate(notEmpty),
			huh.NewInput().Title("Password:").
				Value(&c.password).
				EchoMode(huh.EchoModePassword),
			huh.NewConfirm().Title("Default server settings?").
				Key("defaultServer").
				Value(&c.defaultServer).
				Affirmative("Y").
				Negative("N"),
		),
		huh.NewGroup(
			c.ipInput,
			c.portInput).
			WithHideFunc(func() bool {
				return c.defaultServer
			}),
		// huh.NewGroup(
		// 	huh.NewConfirm().Title("Login?").
		// 		Key("confirm").
		// 		Affirmative("Y").
		// 		Negative("N").
		// 		Value(&c.confirm),
		// ),
	).
		WithLayout(huh.LayoutStack).
		WithShowHelp(false).
		WithShowErrors(false)
	// WithWidth(width).WithHeight(height)

	c.stableWidth = lipgloss.Width(c.form.View())
	c.stableHeight = lipgloss.Height(c.form.View())
	return c.form.Init()
}

// Update implements PopupModel.
func (c *Connect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
	}

	switch c.form.State {
	case huh.StateCompleted:
		return c, c.ConfirmCmd
	case huh.StateAborted:
		return c, c.CancelCmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if ok && key.Matches(keyMsg, config.Keys.Back) {
		return c, c.CancelCmd
	}

	return c, cmd
}

// GetContentWidth implements PopupModel.
func (c *Connect) GetContentWidth() int {
	return c.stableWidth
}

// GetContentHeight implements PopupModel.
func (c *Connect) GetContentHeight() int {
	return c.stableHeight
}

// GetContentSize implements PopupModel.
func (c *Connect) GetContentSize() (int, int) {
	return c.GetContentWidth(), c.GetContentHeight()
}

func (c *Connect) SetMaxSize(width, height int) {
	w := min(width-ui.PopupStyle.GetHorizontalFrameSize(), c.GetContentWidth())
	h := min(height-ui.PopupStyle.GetVerticalFrameSize(), c.GetContentHeight())

	c.form = c.form.WithWidth(w).WithHeight(h)
}

// View implements PopupModel.
func (c *Connect) View() string {
	if c.form == nil || c.form.State != huh.StateNormal {
		return ""
	}

	errors := c.form.Errors()
	errorText := ""
	for _, err := range errors {
		errorText += err.Error() + "\n"
	}

	errorRender := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Width(c.GetContentWidth()).Height(2).
		Render(errorText)

	render := lipgloss.JoinVertical(lipgloss.Left, c.form.View(), errorRender)

	return render
}

func (c *Connect) ConfirmCmd() tea.Msg {
	var ip string
	var port int
	var err error

	if c.defaultServer {
		ip = "localhost"
		port = 0
	} else {
		ip = c.ip
		port, err = strconv.Atoi(c.port)
		if err != nil {
			port = 0
		}
	}

	return widget.ConnectPopupConfirmMsg{
		DriverName:    c.driverName,
		DefaultServer: c.defaultServer,
		Ip:            ip,
		Port:          port,
		Username:      c.username,
		Password:      c.password,
	}
}
