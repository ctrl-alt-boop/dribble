package io

import (
	tea "github.com/charmbracelet/bubbletea"
)

type GoolDbError struct {
	Err error
}

func NewGoolDbError(err error) tea.Cmd {
	return func() tea.Msg {
		return GoolDbError{
			Err: err,
		}
	}
}

func (e GoolDbError) Error() string {
	return e.Err.Error()
}
