package io

import (
	tea "github.com/charmbracelet/bubbletea"
)

type DribbleError struct {
	Err error
}

func NewDribbleError(err error) tea.Cmd {
	return func() tea.Msg {
		return DribbleError{
			Err: err,
		}
	}
}

func (e DribbleError) Error() string {
	return e.Err.Error()
}
