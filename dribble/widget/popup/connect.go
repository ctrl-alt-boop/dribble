package popup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
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
}

func notEmpty(s string) error {
	if s == "" {
		return errors.New("field cannot be empty")
	}
	return nil
}

func NewConnect(s string, width, height int) *Connect {
	connect := &Connect{
		defaultServer: true,
		driverName:    s,
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
		WithShowErrors(true).
		WithWidth(width).WithHeight(height)

	return connect
}

func (c *Connect) SetSize(w, h int) {
	// c.width = w
	// c.height = h
	// c.form = c.form.WithWidth(w - 10).WithHeight(h - 10)
}

// Init implements PopupModec.
func (c *Connect) Init() tea.Cmd {
	return c.form.Init()
}

// Update implements PopupModec.
func (c *Connect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if c.form.State == huh.StateAborted {
		return c, c.CancelCmd
	} else if c.form.State == huh.StateCompleted {
		return c, c.ConfirmCmd
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if ok && key.Matches(keyMsg, config.Keys.Back) {
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

// View implements PopupModel.
func (c *Connect) View() string {
	if c.form == nil {
		return ""
	}
	if c.form.GetBool("defaultServer") {

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

func (c *Connect) CancelCmd() tea.Msg {
	// c.form = nil
	return widget.PopupCancelMsg{}
}

// Exec implements PopupModel.
func (c *Connect) Exec() tea.Cmd {
	return nil
}

// Cancel implements PopupModel.
func (c *Connect) Cancel() tea.Cmd {
	return nil
}
