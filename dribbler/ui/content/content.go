package content

import "fmt"

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

	Content[T ~[]fmt.Stringer | ~string | ~[]string | ~[][]string | StringTable] interface {
		UpdateSize(width, height int)

		Data() any
		Get() T
	}
)

// Totally unnecessary
func ListToString[L ~[]T, T any, F func(T) string](list L, fn F) []string {
	out := make([]string, len(list))
	for i := range list {
		out[i] = fn(list[i])
	}
	return out
}
