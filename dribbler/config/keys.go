package config

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

const (
	SymbolArrowUp    = "↑"
	SymbolArrowDown  = "↓"
	SymbolArrowLeft  = "←"
	SymbolArrowRight = "→"
	SymbolEnter      = "⏎"
	SymbolBackspace  = "⌫"
	SymbolSpace      = "␣"
	SymbolTab        = "⇥"
	SymbolEscape     = "⎋"
)

const (
	ArrowNav = SymbolArrowLeft + SymbolArrowDown + SymbolArrowUp + SymbolArrowRight
	VimNav   = "hjkl"
	Zoom     = "+/-"
)

const (
	ArrowUp    = "up"
	ArrowDown  = "down"
	ArrowLeft  = "left"
	ArrowRight = "right"
)

type KeyMap struct {
	FullHelpFunc func() [][]key.Binding

	Nav  key.Binding // Used as combined navigation keys for help
	Zoom key.Binding // Used as combined +/- keys for help

	CycleViewNext key.Binding // tab
	CycleViewPrev key.Binding // (shift)+CycleViewNext
	Details       key.Binding // i
	New           key.Binding // n
	NewEmpty      key.Binding // ctrl+New

	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding // enter
	Back   key.Binding // esc

	Increase key.Binding //=
	Decrease key.Binding //-

	Help key.Binding //?
	Quit key.Binding
}

func (keys KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{keys.Help, keys.Quit}
}

func (keys KeyMap) FullHelp() [][]key.Binding {
	if keys.FullHelpFunc != nil {
		return keys.FullHelpFunc()
	}
	return [][]key.Binding{
		{keys.Help},
		{keys.Quit},
		{keys.CycleViewNext},
		{keys.Nav},
		{keys.Select},
		{keys.Back},
	}
}

var Keys = createKeyMap()

func createKeyMap() KeyMap {
	return KeyMap{
		Nav: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(ArrowNav+"/"+VimNav, "navigate"),
		),
		CycleViewNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(SymbolTab, "cycle view"),
		),
		CycleViewPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("S+"+SymbolTab, "cycle view prev"),
		),
		Details: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "toggle details"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new ..."),
		),
		NewEmpty: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("C+n", "new empty"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp(SymbolArrowLeft, "left"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(SymbolArrowDown, "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(SymbolArrowUp, "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(SymbolArrowRight, "right"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(SymbolEnter, "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(SymbolEscape, "back"),
		),
		Zoom: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(Zoom, "zoom"),
		),
		Increase: key.NewBinding(
			key.WithKeys("="),
			key.WithHelp("+", "inc"),
		),
		Decrease: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "dec"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

var LoginKeyMap = &huh.KeyMap{
	Input: huh.InputKeyMap{
		AcceptSuggestion: key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "complete")),
		Prev:             key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
		Next:             key.NewBinding(key.WithKeys("enter", "tab"), key.WithHelp("enter", "next")),
		Submit:           key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	},
}
