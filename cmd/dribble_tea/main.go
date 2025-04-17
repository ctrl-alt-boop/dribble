package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

func main() {
	var ip string = "localhost"
	if len(os.Args) > 1 {
		ip = os.Args[1]
	}

	logger := logging.NewLogger("app.log")
	defer logger.Close()

	logger.Info("GoolDb Create")
	gool := gooldb.New(logger, ip)

	dribble := dribble.NewModel(gool)
	p := tea.NewProgram(dribble, tea.WithAltScreen())
	dribble.SetProgramSend(p.Send)
	if _, err := p.Run(); err != nil {
		fmt.Printf("GoolTea error: %v\n", err)
		os.Exit(1)
	}
}
