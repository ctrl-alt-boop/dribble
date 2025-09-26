package content

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var _ Content[[]string] = (*List)(nil)
var _ Selection = (*List)(nil)
var _ tea.Model = (*List)(nil)

type Item struct {
	ID    int
	Index int

	Value any
}

func (l *Item) String() string {
	return fmt.Sprint(l.Value)
}

type List struct {
	Items []Item

	Width, Height int
	cursorY       int
}

func (l *List) GetSelected() any {
	return l.Items[l.cursorY]
}

func (l *List) Cursor() (int, int) {
	return 0, l.cursorY
}

func (l *List) CursorX() int {
	return 0
}

func (l *List) CursorY() int {
	return l.cursorY
}

func (l *List) MoveCursor(_ int, dY int) {
	l.SetCursor(0, l.cursorY+dY)
}

func (l *List) MoveCursorDown(y ...int) {
	l.MoveCursor(0, 1)
}

func (l *List) MoveCursorLeft(_ ...int) {
	// Not applicable for a simple list
}

func (l *List) MoveCursorRight(_ ...int) {
	// Not applicable for a simple list
}

func (l *List) MoveCursorUp(y ...int) {
	l.MoveCursor(0, -1)
}

func (l *List) SetCursor(_ int, y int) {
	if y < 0 {
		y = 0
	}
	if y >= len(l.Items) {
		y = len(l.Items) - 1
	}
	l.cursorY = y
}

func (l *List) Init() tea.Cmd {
	return nil
}

func (l *List) Data() any {
	return l.Items
}

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l *List) UpdateSize(width int, height int) {
	l.Width, l.Height = width, height
}

func (l *List) String() string {
	itemStrings := ListToString(l.Items, func(item Item) string {
		return item.String()
	})
	return strings.Join(itemStrings, "\n")
}

func (l *List) Get() []string {
	s := make([]string, len(l.Items))
	for i, item := range l.Items {
		s[i] = item.String()
	}
	return s
}

func (l *List) View() string {
	return l.String()
}
