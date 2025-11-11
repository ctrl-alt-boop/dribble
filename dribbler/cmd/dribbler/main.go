// Package main
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribbler"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

func main() {
	dribbler := dribbler.NewModel()
	p := tea.NewProgram(dribbler)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Dribble error: %v\n", err)
		os.Exit(1)
	}
	logging.CloseGlobalLogger()
	os.Exit(0)
}
