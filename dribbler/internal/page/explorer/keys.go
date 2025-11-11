package explorer

import (
	"charm.land/bubbles/v2/key"
	"github.com/ctrl-alt-boop/dribbler/internal/keys"
)

type (
	NavigationKeys struct {
		Nav key.Binding

		Left  key.Binding
		Down  key.Binding
		Up    key.Binding
		Right key.Binding
	}

	NumberedSelectKeys struct {
		Number key.Binding
	}

	CursorKeys struct {
		*NavigationKeys

		Accept  key.Binding
		Cancel  key.Binding
		Details key.Binding
		Esc     key.Binding
	}

	ExplorerPageKeys struct {
		FocusSidebar    key.Binding // F1, C+1
		FocusWorkspace  key.Binding // F2, C+2
		CommandlineMode key.Binding // :

		FocusCycleNext key.Binding // tab
		FocusCyclePrev key.Binding // shift+tab

		ToggleHelp key.Binding // ?
	}

	SidebarKeys struct {
		*CursorKeys
		*NumberedSelectKeys

		CreateNew key.Binding
		Search    key.Binding
	}

	WorkspaceKeys struct {
		*CursorKeys
		*NumberedSelectKeys
		NextTab key.Binding
		PrevTab key.Binding
	}

	CommandlineKeys struct {
		Submit key.Binding
		Cancel key.Binding
	}

	CreationPaneKeys struct {
		*NumberedSelectKeys
		Accept key.Binding
		Cancel key.Binding

		Query key.Binding
		//...
	}
)

func DefaultNavigationKeys() *NavigationKeys {
	return &NavigationKeys{
		// Navigation
		Nav: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(keys.ArrowNav+"/"+keys.VimNav, "cursor"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp(keys.SymbolArrowLeft, "left"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(keys.SymbolArrowDown, "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(keys.SymbolArrowUp, "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(keys.SymbolArrowRight, "right"),
		),
	}
}

func DefaultNumberedSelectKeys() *NumberedSelectKeys {
	return &NumberedSelectKeys{
		Number: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9", "0"),
			key.WithHelp("0-9", "select"),
		),
	}
}

func DefaultCursorKeys() *CursorKeys {
	return &CursorKeys{
		NavigationKeys: DefaultNavigationKeys(),
		Accept: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(keys.SymbolEnter, "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp(keys.SymbolBackspace, "back"),
		),
		Details: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "toggle details"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(keys.SymbolEscape, "escape"),
		),
	}
}

func DefaultExplorerKeyBindings() *ExplorerPageKeys {
	return &ExplorerPageKeys{
		// Focus
		FocusSidebar: key.NewBinding(
			key.WithKeys("f1", "ctrl+1"),
			key.WithHelp("C+1", "side panel"),
		),
		FocusWorkspace: key.NewBinding(
			key.WithKeys("f2", "ctrl+2"),
			key.WithHelp("C+2", "workspace"),
		),
		CommandlineMode: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "command-line mode"),
		),

		FocusCycleNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(keys.SymbolTab, "cycle focus next"),
		),
		FocusCyclePrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("S+"+keys.SymbolTab, "cycle focus prev"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

func DefaultSidebarKeyBindings() *SidebarKeys {
	return &SidebarKeys{
		CursorKeys:         DefaultCursorKeys(),
		NumberedSelectKeys: DefaultNumberedSelectKeys(),
		CreateNew: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new ..."),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
	}
}

func DefaultWorkspaceKeyBindings() *WorkspaceKeys {
	return &WorkspaceKeys{
		CursorKeys:         DefaultCursorKeys(),
		NumberedSelectKeys: DefaultNumberedSelectKeys(),
		NextTab: key.NewBinding(
			key.WithKeys("]"),
			key.WithHelp("]", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("["),
			key.WithHelp("[", "prev tab"),
		),
	}
}

func DefaultCommandlineKeyBindings() *CommandlineKeys {
	return &CommandlineKeys{
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(keys.SymbolEnter, "submit"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(keys.SymbolEscape, "cancel"),
		),
	}
}

func DefaultCreationPaneKeyBindings() *CreationPaneKeys {
	return &CreationPaneKeys{
		NumberedSelectKeys: DefaultNumberedSelectKeys(),
		Cancel: key.NewBinding(
			key.WithKeys("esc", "backspace"),
			key.WithHelp("🗙", "escape"),
		),
	}
}

// ShortHelp implements help.KeyMap
func (n NavigationKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		n.Nav,
	}
}

// FullHelp implements help.KeyMap
func (n NavigationKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{n.Nav},
		{n.Left},
		{n.Down},
		{n.Up},
		{n.Right},
	}
}

// ShortHelp implements help.KeyMap
func (n NumberedSelectKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		n.Number,
	}
}

// FullHelp implements help.KeyMap
func (n NumberedSelectKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{n.Number},
	}
}

// ShortHelp implements help.KeyMap
func (c CursorKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		c.NavigationKeys.Nav,
		c.Accept,
		c.Cancel,
		c.Details,
	}
}

// FullHelp implements help.KeyMap
func (c CursorKeys) FullHelp() [][]key.Binding {
	return append(
		c.NavigationKeys.FullHelp(),
		[][]key.Binding{
			{c.Accept},
			{c.Cancel},
			{c.Details},
			{c.Esc},
		}...,
	)
}

// ShortHelp implements help.KeyMap
func (e ExplorerPageKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		e.FocusSidebar,
		e.FocusWorkspace,
		e.FocusCycleNext,
		e.ToggleHelp,
	}
}

// FullHelp implements help.KeyMap
func (e ExplorerPageKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{e.FocusSidebar},
		{e.FocusWorkspace},
		{e.CommandlineMode},
		{e.FocusCycleNext},
		{e.FocusCyclePrev},
		{e.ToggleHelp},
	}
}

// ShortHelp implements help.KeyMap
func (s SidebarKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		s.CursorKeys.Nav,
		s.NumberedSelectKeys.Number,
		s.CreateNew,
	}
}

// FullHelp implements help.KeyMap
func (s SidebarKeys) FullHelp() [][]key.Binding {
	k := []key.Binding{s.CreateNew}
	return append(
		append(
			s.CursorKeys.FullHelp(),
			s.NumberedSelectKeys.FullHelp()...,
		),
		k,
	)
}

// ShortHelp implements help.KeyMap
func (w WorkspaceKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		w.CursorKeys.Nav,
		w.NumberedSelectKeys.Number,
	}
}

// FullHelp implements help.KeyMap
func (w WorkspaceKeys) FullHelp() [][]key.Binding {
	return append(
		w.CursorKeys.FullHelp(),
		w.NumberedSelectKeys.FullHelp()...,
	)
}

// ShortHelp implements help.KeyMap
func (c CommandlineKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		c.Submit,
		c.Cancel,
	}
}

// FullHelp implements help.KeyMap
func (c CommandlineKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{c.Submit},
		{c.Cancel},
	}
}

// ShortHelp implements help.KeyMap
func (c CreationPaneKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		c.NumberedSelectKeys.Number,
		c.Cancel,
	}
}

// FullHelp implements help.KeyMap
func (c CreationPaneKeys) FullHelp() [][]key.Binding {
	return append(
		c.NumberedSelectKeys.FullHelp(),
		[]key.Binding{c.Cancel},
	)
}
