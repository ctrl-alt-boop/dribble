package content

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

var logger = logging.GlobalLogger()

type (
	Selection interface {
		CursorX() int
		CursorY() int
		Cursor() (int, int)

		SetCursor(x, y int)
		MoveCursor(dX, dY int)

		MoveCursorUp(y ...int)
		MoveCursorDown(y ...int)
		MoveCursorLeft(x ...int)
		MoveCursorRight(x ...int)

		GetSelected() any
	}
)

var DefaultStyle = lipgloss.NewStyle().Margin(1, 2)

// Totally unnecessary
func ListToString[L ~[]T, T any, F func(T) string](list L, fn F) []string {
	out := make([]string, len(list))
	for i := range list {
		out[i] = fn(list[i])
	}
	return out
}
