package popup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
)

type Connect struct {
	form      *huh.Form
	ipInput   *huh.Input
	portInput *huh.Input

	confirm bool

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
	formTitle := fmt.Sprintf("Connect to %s server", s)
	connect.ipInput = huh.NewInput().
		Key("ip").
		Value(&connect.ip).
		Title("IP:")
	connect.portInput = huh.NewInput().
		Key("port").
		Value(&connect.port).
		Title("Port:")

	connect.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(formTitle),

			huh.NewInput().
				Value(&connect.username).
				Validate(notEmpty).
				Title("Username:"),
			huh.NewInput().
				Value(&connect.password).
				EchoMode(huh.EchoModePassword).
				Title("Password:"),

			huh.NewConfirm().Title("Default server settings?").
				Key("defaultServer").
				Value(&connect.defaultServer).
				Affirmative("Y").
				Negative("N"),
			connect.ipInput,
			connect.portInput,
			huh.NewConfirm().
				Title("Login?").
				Affirmative("Y").
				Negative("N").
				Value(&connect.confirm),
		),
	).
		WithLayout(huh.LayoutStack).
		WithShowHelp(false).
		WithShowErrors(true)
		// WithWidth(width).WithHeight(height)

	connect.stableWidth = lipgloss.Width(connect.form.View())
	connect.stableHeight = lipgloss.Height(connect.form.View())

	return connect
}

// Init implements PopupModel.
func (c *Connect) Init() tea.Cmd {
	return c.form.Init()
}

// Update implements PopupModel.
func (c *Connect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if c.form.State == huh.StateAborted {
		return c, c.CancelCmd
	} else if c.form.State == huh.StateCompleted {
		return c, c.ConfirmCmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if ok && key.Matches(keyMsg, config.Keys.Back, config.Keys.Quit) {
		return c, c.CancelCmd
	}

	if _, ok := msg.(tea.WindowSizeMsg); ok {
		return c, nil
	}

	form, cmd := c.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		c.form = f
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
	if c.form == nil {
		return ""
	}
	if c.form.GetBool("defaultServer") {
		c.ipInput.Blur()
		c.portInput.Blur()
	} else {
	}

	return c.form.View()
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

	return widget.PopupConfirmMsg{
		DriverName:    c.driverName,
		DefaultServer: c.defaultServer,
		Ip:            ip,
		Port:          port,
		Username:      c.username,
		Password:      c.password,
	}
}

// Exec implements PopupModel.
func (c *Connect) Exec() tea.Cmd {
	return nil
}

// Cancel implements PopupModel.
func (c *Connect) Cancel() tea.Cmd {
	return nil
}
