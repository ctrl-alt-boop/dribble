package explorer

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/ctrl-alt-boop/dribbler/keys"
)

type KeyMap struct {
	FocusSidePanel key.Binding
	FocusWorkspace key.Binding
	FocusPrompt    key.Binding

	FocusCycleNext key.Binding // tab
	FocusFocusPrev key.Binding // (shift)+CycleFocusNext

	Select  key.Binding // enter
	Back    key.Binding // backspace
	Esc     key.Binding // esc
	Details key.Binding // i
	Number  key.Binding // keys 1, 2, ..., 9, 0

	Nav   key.Binding // Used as combined navigation keys for help
	Left  key.Binding
	Down  key.Binding
	Up    key.Binding
	Right key.Binding

	NewQuery key.Binding
}

var keyMap = DefaultKeyMap()

func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Focus
		FocusSidePanel: key.NewBinding(
			key.WithKeys("ctrl+1"),
			key.WithHelp("C+1", "side panel"),
		),
		FocusWorkspace: key.NewBinding(
			key.WithKeys("ctrl+2"),
			key.WithHelp("C+2", "workspace"),
		),
		FocusPrompt: key.NewBinding(
			key.WithKeys("ctrl+3"),
			key.WithHelp("C+3", "prompt"),
		),

		FocusCycleNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp(keys.SymbolTab, "cycle focus next"),
		),
		FocusFocusPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("S+"+keys.SymbolTab, "cycle focus prev"),
		),

		// Interact
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp(keys.SymbolEnter, "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp(keys.SymbolBackspace, "back"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp(keys.SymbolEscape, "escape"),
		),
		Details: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "toggle details"),
		),
		Number: key.NewBinding( // FIXME: I would like to implement multi-numbered select
			key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
			key.WithHelp("1-9", "select"),
		),

		// Navigation
		Nav: key.NewBinding(
			key.WithKeys("nil"),
			key.WithHelp(keys.ArrowNav+"/"+keys.VimNav, "navigate"),
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

		// Creation
		NewQuery: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new ..."),
		),
	}
}
