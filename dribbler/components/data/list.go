package data

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

type Item struct {
	ID    int
	Index int
	Style lipgloss.Style
	Value any
}

func (l *Item) String() string {
	return fmt.Sprintf("%v", l.Value)
}

type List struct {
	Base
	Items []Item

	cursorY int
}

func NewList(items ...Item) List {
	return List{
		Items: items,
	}
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

func (l *List) Set(items ...Item) {
	l.Items = items
}

func (l *List) String() string {
	itemStrings := ListToString(l.Items, func(item Item) string {
		return item.String()
	})
	return strings.Join(itemStrings, "\n")
}

func (l List) Render() string {
	// style := DefaultStyle

	return l.String()
}
