package data_test

import (
	"testing"

	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/ctrl-alt-boop/dribbler/components/data"
)

func newTestResultTable() *result.Table {
	return result.NewTable(
		[]*result.Column{
			{Name: "ID", DBType: "int"},
			{Name: "Name", DBType: "string"},
			{Name: "Value", DBType: "string"},
		},
		[]*result.Row{
			{Values: []any{1, "First", "100"}},
			{Values: []any{2, "Second", "200"}},
		},
	)
}

func TestTable_NewTable(t *testing.T) {
	resTable := newTestResultTable()
	table := data.NewTable(resTable)

	if table == nil {
		t.Fatal("NewTable returned nil")
	}

	if table.Name != "Test Table" {
		t.Errorf("Expected table name 'Test Table', got '%s'", table.Name)
	}
}

func TestTable_CursorMovement(t *testing.T) {
	resTable := newTestResultTable()
	table := data.NewTable(resTable) // 2 rows, 3 cols

	testCases := []struct {
		name         string
		action       func()
		expectedX    int
		expectedY    int
		setup        func()
		moveX, moveY int
	}{
		{name: "SetCursor inside bounds", setup: func() { table.SetCursor(1, 1) }, expectedX: 1, expectedY: 1},
		{name: "SetCursor outside bounds (negative)", setup: func() { table.SetCursor(-1, -5) }, expectedX: 0, expectedY: 0},
		{name: "SetCursor outside bounds (positive)", setup: func() { table.SetCursor(10, 10) }, expectedX: 2, expectedY: 1},
		{name: "MoveCursorRight", setup: func() { table.SetCursor(0, 0) }, action: func() { table.MoveCursorRight() }, expectedX: 1, expectedY: 0},
		{name: "MoveCursorLeft from edge", setup: func() { table.SetCursor(0, 0) }, action: func() { table.MoveCursorLeft() }, expectedX: 0, expectedY: 0},
		{name: "MoveCursorDown", setup: func() { table.SetCursor(0, 0) }, action: func() { table.MoveCursorDown() }, expectedX: 0, expectedY: 1},
		{name: "MoveCursorUp from edge", setup: func() { table.SetCursor(0, 0) }, action: func() { table.MoveCursorUp() }, expectedX: 0, expectedY: 0},
		{name: "MoveCursor generic", setup: func() { table.SetCursor(0, 0) }, action: func() { table.MoveCursor(1, 1) }, expectedX: 1, expectedY: 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			if tc.action != nil {
				tc.action()
			}
			x, y := table.Cursor()
			if x != tc.expectedX || y != tc.expectedY {
				t.Errorf("Expected cursor at (%d, %d), but got (%d, %d)", tc.expectedX, tc.expectedY, x, y)
			}
		})
	}
}

func TestTable_GetSelected(t *testing.T) {
	resTable := newTestResultTable()
	table := data.NewTable(resTable)

	table.SetCursor(1, 1) // "Second"
	selectedValue := table.GetSelected()

	if val, ok := selectedValue.(string); !ok || val != "Second" {
		t.Errorf("Expected selected value to be 'Second', but got '%v'", selectedValue)
	}

	table.SetCursor(2, 0) // "100"
	selectedValue = table.GetSelected()

	if val, ok := selectedValue.(string); !ok || val != "100" {
		t.Errorf("Expected selected value to be '100', but got '%v'", selectedValue)
	}
}

func TestTable_Update(t *testing.T) {
	resTable := newTestResultTable()
	table := data.NewTable(resTable)

	// Create new data for the update
	newCols := []*result.Column{{Name: "Status", DBType: "string"}}
	newRows := []*result.Row{
		{Values: []any{"Pending"}},
		{Values: []any{"OK"}},
		{Values: []any{"Done"}},
	}
	newResTable := result.NewTable(newCols, newRows)

	// Send the correct update
	table.Set(newResTable)
}
