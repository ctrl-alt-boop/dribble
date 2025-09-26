package content

import (
	"fmt"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/result"
)

const DefaultMaxCellWidth int = 36 // Guid length, including the '-'s

type StringTable0 struct {
	Columns []string
	Rows    [][]string
	Widths  []int
}

type StringTable struct {
	// map[Column]ColumnRows
	Data   map[string][]string
	Widths []int
}

func (s StringTable) String() string {
	return ""
}

var _ Content[StringTable] = (*Table)(nil)
var _ Selection = (*Table)(nil)
var _ tea.Model = (*Table)(nil)

type SetTableContentMsg struct {
	ID    int
	Name  string
	Table *result.Table
}

var _ fmt.Stringer = (*Cell)(nil)

type Cell struct {
	ID                   int
	LocationX, LocationY int
	Value                any
	Placeholder          string // Maybe
}

// String implements fmt.Stringer.
func (c *Cell) String() string {
	return fmt.Sprint(c.Value)
}

type Table struct {
	ID    int
	Name  string
	Table *result.Table

	MaxCellWidth  int
	Width, Height int
	columnWidths  []int

	cursorX int
	cursorY int

	rowTextTemplate *template.Template
}

// GetSelected implements Cursored.
func (t *Table) GetSelected() any {
	value, err := t.Table.GetRowColumn(t.cursorY, t.cursorX)
	if err != nil {
		return "VALUE ERROR"
	}
	return value
}

// Cursor implements Cursored.
func (t *Table) Cursor() (int, int) {
	return t.cursorX, t.cursorY
}

// CursorX implements Cursored.
func (t *Table) CursorX() int {
	return t.cursorX
}

// CursorY implements Cursored.
func (t *Table) CursorY() int {
	return t.cursorY
}

// MoveCursor implements Content.
func (t *Table) MoveCursor(dX int, dY int) {
	t.SetCursor(t.cursorX+dX, t.cursorY+dY)
}

// MoveCursorDown implements Content.
func (t *Table) MoveCursorDown(y ...int) {
	t.MoveCursor(0, 1)
}

// MoveCursorLeft implements Content.
func (t *Table) MoveCursorLeft(x ...int) {
	t.MoveCursor(-1, 0)
}

// MoveCursorRight implements Content.
func (t *Table) MoveCursorRight(x ...int) {
	t.MoveCursor(1, 0)
}

// MoveCursorUp implements Content.
func (t *Table) MoveCursorUp(y ...int) {
	t.MoveCursor(0, -1)
}

// SetCursor implements Content.
func (t *Table) SetCursor(x int, y int) {
	if t.Table == nil {
		return
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= t.Table.NumColumns() {
		x = t.Table.NumColumns() - 1
	}
	if y >= t.Table.NumRows() {
		y = t.Table.NumRows() - 1
	}
	t.cursorX = x
	t.cursorY = y
}

func NewTable(id int, name string, table *result.Table) *Table {
	newTable := &Table{
		ID:    id,
		Name:  name,
		Table: table,

		MaxCellWidth: DefaultMaxCellWidth,
	}
	newTable.columnWidths = getColumnWidths(table, newTable.MaxCellWidth)
	return newTable
}

// Init implements tea.Model.
func (t *Table) Init() tea.Cmd {
	return nil
}

// Data implements Content.
func (t *Table) Data() any {
	return t.Table
}

// Update implements Content.
func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetTableContentMsg:
		if msg.ID != t.ID {
			return t, nil
		}
		t.Name = msg.Name
		t.Table = msg.Table
		t.columnWidths = getColumnWidths(msg.Table, t.MaxCellWidth)
		t.rowTextTemplate = template.Must(
			template.New("row").
				Funcs(template.FuncMap{
					"fixLength": t.fixLength,
				}).
				Parse(rowTextTemplate))
		return t, nil
	}
	return t, nil
}

// UpdateSize implements Content.
func (t *Table) UpdateSize(width int, height int) {
	t.Width, t.Height = width, height
}

const rowTextTemplate = "\u2502 {{- range $i, $e := . }} {{fixLength $e $i}} \u2502{{- end -}}"

func truncPad(s string, min, max int) string {
	if len(s) >= max {
		return s[:max]
	}
	if len(s) <= min {
		return fmt.Sprintf("%-*s", min, s)
	}
	return s
}

func (t *Table) fixLength(s string, columnIndex int) string {
	return truncPad(s, t.columnWidths[columnIndex], t.MaxCellWidth)
}

func (t *Table) GetV2() StringTable0 {
	if t.Table == nil {
		return StringTable0{}
	}
	return StringTable0{
		Columns: t.Table.ColumnNames(),
		Rows:    t.Table.GetRowStringsAll(),
	}
}

func (t *Table) Get() StringTable {
	if t.Table == nil {
		return StringTable{}
	}
	data := make(map[string][]string, t.Table.NumColumns())
	widths := make([]int, t.Table.NumColumns())
	for i := range t.Table.Columns() {
		rows, width := t.Table.GetColumnRows(i)
		data[t.Table.Columns()[i].Name] = rows
		widths[i] = width
	}
	return StringTable{
		Data:   data,
		Widths: widths,
	}
}

func (t *Table) View() string {
	if t.Table == nil {
		return ""
	}
	var sb strings.Builder

	columns := t.Table.ColumnNames()
	sb.WriteString(strings.Join(columns, "\t"))
	sb.WriteString("\n")

	for _, row := range t.Table.GetRowStringsAll() {
		err := t.rowTextTemplate.Execute(&sb, row)
		if err != nil {
			return err.Error()
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func getColumnWidths(dataTable *result.Table, maxCellWidth int) []int {
	columnWidths := make([]int, dataTable.NumColumns())
	for i := range dataTable.Rows() {
		row := dataTable.GetRowStrings(i)
		for columnIndex, value := range row {
			if len(value) >= maxCellWidth {
				columnWidths[columnIndex] = maxCellWidth
			}
			if len(value) > columnWidths[columnIndex] && len(value) <= maxCellWidth {
				columnWidths[columnIndex] = len(value)
			}
			if len(dataTable.Columns()[columnIndex].Name) >= columnWidths[columnIndex] {
				columnWidths[columnIndex] = len(dataTable.Columns()[columnIndex].Name) + 2
			}
		}
	}
	return columnWidths
}
