package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler"
)

func main() {

	dribble := dribble.NewClient()

	dribbler := dribbler.CreateDemoModel(dribble)
	p := tea.NewProgram(dribbler, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Dribble error: %v\n", err)
		os.Exit(1)
	}
}
