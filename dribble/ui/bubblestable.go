package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/ctrl-alt-boop/gooldb/dribble/util"
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
)

const defaultCellWidth = 36 // Guid length, including the '-'s

type BubblesTable struct {
	Table table.Model

	maxCellWidth int
}

func (t *BubblesTable) View() string {
	return t.Table.View()
}

func New() *BubblesTable {
	return &BubblesTable{
		maxCellWidth: defaultCellWidth,
	}
}

func (t *BubblesTable) SetTable(dataTable data.Table) {
	newTable := table.New()
	tableColumns, widths := t.getColumnsAsTableColumns(dataTable)
	newTable.SetColumns(tableColumns)

	tableRows := t.getRowsAsTableRows(dataTable)
	newTable.SetRows(tableRows)

	t.Table.SetHeight(len(tableRows) + 3)
	t.Table.SetWidth(util.Sum(widths...))
	// t.Table.SetStyles()

	t.Table = newTable
}

func (t *BubblesTable) IsTableSet() bool {
	return len(t.Table.Columns()) > 0
}

func (t *BubblesTable) getColumnsAsTableColumns(dataTable data.Table) ([]table.Column, []int) {
	columnNames := dataTable.ColumnNames()
	columnWidths := t.GetColumnWidths(dataTable)
	columns := util.Zip(columnNames, columnWidths)
	var tableColumns []table.Column
	for _, column := range columns {
		tableColumns = append(tableColumns, table.Column{Title: column.First, Width: column.Second})
	}
	return tableColumns, columnWidths
}

func (t *BubblesTable) getRowsAsTableRows(dataTable data.Table) []table.Row {
	var tableRows []table.Row
	for i := range dataTable.Rows() {
		row := dataTable.GetRowStrings(i)
		tableRows = append(tableRows, table.Row(row))
	}
	return tableRows
}

func (t *BubblesTable) GetColumnWidths(dataTable data.Table) []int {
	columnWidths := make([]int, dataTable.NumColumns())
	for i := range dataTable.Rows() {
		row := dataTable.GetRowStrings(i)
		for columnIndex, value := range row {
			if len(value) >= t.maxCellWidth {
				columnWidths[columnIndex] = t.maxCellWidth
			}
			if len(value) > columnWidths[columnIndex] && len(value) <= t.maxCellWidth {
				columnWidths[columnIndex] = len(value)
			}
			if len(dataTable.Columns()[columnIndex].Name) > columnWidths[columnIndex] {
				columnWidths[columnIndex] = len(dataTable.Columns()[columnIndex].Name)
			}
		}
	}
	return columnWidths
}
