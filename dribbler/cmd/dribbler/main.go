// Package main
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribbler"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

func main() {
	dribbler := dribbler.NewDribblerModel()
	p := tea.NewProgram(dribbler, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Dribble error: %v\n", err)
		os.Exit(1)
	}
	logging.CloseGlobalLogger()
	os.Exit(0)
}
