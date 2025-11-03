package data_test

import (
	"fmt"
	"testing"

	"github.com/ctrl-alt-boop/dribbler/components/data"
)

// const rowTextTemplate = "\u2502 {{- range $i, $e := . }} {{fixLength $e $i}} \u2502{{- end -}}"

// func truncPad(s string, min, max int) string {
// 	if len(s) >= max {
// 		return s[:max]
// 	}
// 	if len(s) <= min {
// 		return fmt.Sprintf("%-*s", min, s)
// 	}
// 	return s
// }

// func TestRowTemplate(t *testing.T) {
// 	rows := []string{"foooooo", "bar", "bazzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"}

// 	columnWidths := []int{6, 6, 62}
// 	maxCellWidth := 8

// 	fixLength := func(s string, columnIndex int) string {
// 		fmt.Println("fixLength", s, " ", columnIndex, " ", columnWidths[columnIndex], " ", maxCellWidth)
// 		return truncPad(s, columnWidths[columnIndex], maxCellWidth)
// 	}

// 	tmpl := template.Must(template.New("row").Funcs(template.FuncMap{
// 		"fixLength": fixLength,
// 	}).Parse(rowTextTemplate))

// 	var sb strings.Builder
// 	err := tmpl.Execute(&sb, rows)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	res := sb.String()
// 	if res != "│ foooooo │ bar    │ bazzzzzz │" {
// 		runes := []rune(res)
// 		for r := range runes {
// 			fmt.Printf("'%v'", string(runes[r]))
// 		}
// 		fmt.Printf("\n")

//			runes = []rune("│ foooooo │ bar    │ bazzzzzz │")
//			for r := range runes {
//				fmt.Printf("'%v'", string(runes[r]))
//			}
//			fmt.Printf("\n")
//			fmt.Println(res)
//			fmt.Println("│ foooooo │ bar    │ bazzzzzz │")
//			t.Fail()
//		}
//	}

func TestListToString(t *testing.T) {
	list := []data.Item{
		{Value: "foo"},
		{Value: "bar"},
		{Value: "baz"},
		{Value: "qux"},
	}
	s := data.ListToString(list, func(item data.Item) string {
		return fmt.Sprint(item.Value)
	})

	fmt.Print(s)

	// table := [][]Item{}
	// l := ToString[[]Item, Item, func(Item) string](table, func(item Item) string {
	// 	s := make([]string, len(table))
	// 	for i, item := range table {
	// 		s[i] = fmt.Sprint(item[0])
	// 	}
	// 	return s
	// })

	// tree := []*Node{}

	// t := ToString(tree, func(node *Node) string {
	// 	return node
	// })
}
