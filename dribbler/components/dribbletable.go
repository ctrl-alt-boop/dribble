package components

import (
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/ctrl-alt-boop/dribbler/util"
)

const defaultCellWidth = 36 // Guid length, including the '-'s

var (
	mainStyle = lipgloss.NewStyle().
			Margin(0, 1).
			Height(1).
			AlignHorizontal(lipgloss.Center)

	cellStyle = mainStyle.
			Faint(true)

	headerStyle = lipgloss.NewStyle().
			Padding(0, 2)

	highlightedRowStyle = mainStyle

	highlightedCellStyle = mainStyle.
				Background(lipgloss.Color("7"))

	headerBoxStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true)
)

type (
	Row    []string
	Rows   []Row
	Column struct {
		Value string
		Width int
	}
	Columns []Column
)

type DribbleTable struct {
	Height, Width int

	ViewportWidth, ViewportHeight int

	Columns Columns
	Rows    Rows

	selection struct {
		row, column int
	}
	Offset struct {
		X, Y int
	}

	maxColumnWidth int
}

func (t *DribbleTable) IncreaseColumnSize() {
	t.Columns[t.selection.column].Width++
}

func (t *DribbleTable) DecreaseColumnSize() {
	t.Columns[t.selection.column].Width--
}

func (t *DribbleTable) GetSelected() string {
	return t.Rows[t.selection.row][t.selection.column]
}

func (t *DribbleTable) ColumnWidths() []int {
	widths := make([]int, len(t.Columns))
	for i, column := range t.Columns {
		widths[i] = column.Width
	}
	return widths
}

func (t *DribbleTable) MoveCursorUp() {
	t.selection.row--
	if t.selection.row < 0 {
		t.selection.row = 0
	}
}

func (t *DribbleTable) MoveCursorDown() {
	t.selection.row++
	if t.selection.row >= len(t.Rows) {
		t.selection.row = len(t.Rows) - 1
	}
}

func (t *DribbleTable) MoveCursorLeft() {
	t.selection.column--
	if t.selection.column < 0 {
		t.selection.column = 0
	}
	t.scrollToActiveColumn()
}

func (t *DribbleTable) MoveCursorRight() {
	t.selection.column++
	if t.selection.column >= len(t.Columns) {
		t.selection.column = len(t.Columns) - 1
	}
	t.scrollToActiveColumn()
}

func (t *DribbleTable) scrollToActiveColumn() {
	if len(t.ColumnWidths()) == 0 {
		return
	}

	columnWidths := t.ColumnWidths()

	activeColumnStart := 0
	for i := range columnWidths[:t.selection.column] {
		activeColumnStart += columnWidths[i]
	}

	activeColumnEnd := activeColumnStart + columnWidths[t.selection.column]

	edgeLeft := t.Offset.X
	edgeRight := t.Offset.X + t.ViewportWidth

	if activeColumnStart < edgeLeft { // Case 1: Active column starts BEFORE the visible part of the viewport (scroll left)
		t.Offset.X = activeColumnStart
	} else if activeColumnEnd > edgeRight { // Case 2: Active column ends AFTER the visible part of the viewport (scroll right)
		t.Offset.X = activeColumnEnd - t.ViewportWidth + 1 // Just a small 'padding' like value
	}

	if t.Offset.X < 0 {
		t.Offset.X = 0
	}
	if t.Width > t.ViewportWidth {
		if t.Offset.X > t.Width-t.ViewportWidth {
			t.Offset.X = t.Width - t.ViewportWidth
		}
	} else {
		t.Offset.X = 0
	}
}

func (t *DribbleTable) View() string {
	if !t.IsTableSet() {
		return ""
	}
	headerCells := make([]string, len(t.Columns))
	for i, column := range t.Columns {
		headerCells[i] = headerStyle.Width(column.Width).Render(column.Value)
	}
	header := headerBoxStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))

	rows := make([]string, len(t.Rows))
	for i, row := range t.Rows {
		cells := make([]string, len(t.Columns))
		for j, cell := range row {
			cellSize := t.Columns[j].Width - mainStyle.GetHorizontalFrameSize()
			value := util.TruncateWithSuffix(cell, cellSize, "...")
			style := cellStyle
			if i == t.selection.row && j == t.selection.column {
				style = highlightedCellStyle
			} else if i == t.selection.row {
				style = highlightedRowStyle
			}

			cells[j] = style.Width(cellSize).AlignHorizontal(lipgloss.Left).Render(value)
		}

		rows[i] = lipgloss.JoinHorizontal(lipgloss.Top, cells...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func NewDribbleTable() *DribbleTable {
	return &DribbleTable{
		maxColumnWidth: defaultCellWidth,
	}
}

func (t *DribbleTable) SetTable(dataTable result.Table) {
	tableColumns := t.getColumnsAsTableColumns(dataTable)
	t.Columns = tableColumns

	tableRows := t.getRowsAsTableRows(dataTable)
	t.Rows = tableRows

	widths := make([]int, len(tableColumns))
	for i, column := range t.Columns {
		widths[i] = column.Width
	}

	t.Height = len(tableRows) + 3
	t.Width = util.Sum(widths...)

	t.selection.row = 0
	t.selection.column = 0
	t.Offset.X = 0
	t.Offset.Y = 0
}

func (t *DribbleTable) IsTableSet() bool {
	return len(t.Columns) > 0
}

func (t *DribbleTable) IsTableEmpty() bool {
	return len(t.Rows) == 0
}

func (t *DribbleTable) getColumnsAsTableColumns(dataTable result.Table) Columns {
	columnNames := dataTable.ColumnNames()
	columnWidths := t.getColumnWidths(dataTable)
	columns := util.Zip(columnNames, columnWidths)
	var tableColumns Columns
	for _, column := range columns {
		tableColumns = append(tableColumns, Column{Value: column.Left, Width: column.Right + 2})
	}
	return tableColumns
}

func (t *DribbleTable) getRowsAsTableRows(dataTable result.Table) Rows {
	var tableRows Rows
	for i := range dataTable.Rows() {
		row := dataTable.GetRowStrings(i)
		tableRows = append(tableRows, Row(row))
	}
	return tableRows
}

func (t *DribbleTable) getColumnWidths(dataTable result.Table) []int {
	columnWidths := make([]int, dataTable.NumColumns())
	for i := range dataTable.Rows() {
		row := dataTable.GetRowStrings(i)
		for columnIndex, value := range row {
			if len(value) >= t.maxColumnWidth {
				columnWidths[columnIndex] = t.maxColumnWidth
			}
			if len(value) > columnWidths[columnIndex] && len(value) <= t.maxColumnWidth {
				columnWidths[columnIndex] = len(value)
			}
			if len(dataTable.Columns()[columnIndex].Name) >= columnWidths[columnIndex] {
				columnWidths[columnIndex] = len(dataTable.Columns()[columnIndex].Name) + 2
			}
		}
	}
	return columnWidths
}
