package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler"
)

func main() {
	// defer logging.CloseGlobalLogger()
	var ip string = "localhost"
	if len(os.Args) > 1 {
		ip = os.Args[1]
	}

	// logging.GlobalLogger().Infof("DribbleAPI Create")
	dribble := dribble.NewClient(ip)

	dribbleTUI := dribbler.InitialModel(dribble)
	p := tea.NewProgram(dribbleTUI, tea.WithAltScreen())
	dribbleTUI.SetProgramSend(p.Send)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Dribble error: %v\n", err)
		os.Exit(1)
	}
}
