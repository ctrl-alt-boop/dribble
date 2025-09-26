package dribbler

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/request"
)

type TableNamesMapMsg struct {
	TableNames map[string][]string
	Err        error
}

func (c *AppModel) GetAllTableNames() tea.Cmd {
	return func() tea.Msg {

		respChan, err := c.dribbleClient.RequestForAll(context.TODO(), request.NewReadTableNames())
		if err != nil {
			return err
		}
		tableNames := map[string][]string{}
		for resp := range respChan {
			tableNames[resp.RequestTarget] = resp.Body.([]string)
		}
		return TableNamesMapMsg{
			TableNames: tableNames,
		}
	}
}
