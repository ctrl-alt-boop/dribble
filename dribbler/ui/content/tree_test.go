package content_test

import (
	"reflect"
	"testing"

	"github.com/ctrl-alt-boop/dribbler/ui/content"
)

// newTestTree creates a standard tree for testing navigation.
// Structure:
// - root (expanded)
//   - child 0 (expanded)
//   - grandchild 0.0
//   - grandchild 0.1
//   - child 1 (collapsed)
//   - grandchild 1.0
//   - child 2
func newTestTree() *content.Tree {
	tree := content.NewTree() // root is expanded by default
	c0 := tree.NewChild("child 0", nil)
	c0.NewChild("grandchild 0.0", nil)
	c0.NewChild("grandchild 0.1", nil)
	c0.Expanded = true

	c1 := tree.NewChild("child 1", nil)
	c1.NewChild("grandchild 1.0", nil)
	c1.Expanded = false // collapsed

	tree.NewChild("child 2", nil)

	return tree
}

func TestNode_NewChild_And_GetAtPath(t *testing.T) {
	tree := newTestTree()

	// Test GetAtPath
	if node := tree.GetAtPath([]int{0}); node.Name != "child 0" {
		t.Errorf("Expected to get 'child 0' at path [0], but got '%s'", node.Name)
	}
	if node := tree.GetAtPath([]int{0, 1}); node.Name != "grandchild 0.1" {
		t.Errorf("Expected to get 'grandchild 0.1' at path [0, 1], but got '%s'", node.Name)
	}
	if node := tree.GetAtPath([]int{2}); node.Name != "child 2" {
		t.Errorf("Expected to get 'child 2' at path [2], but got '%s'", node.Name)
	}

	// Test IndexPath correctness
	expectedPath := []int{1, 0}
	actualPath := tree.GetAtPath([]int{1, 0}).IndexPath
	if !reflect.DeepEqual(actualPath, expectedPath) {
		t.Errorf("Expected IndexPath for grandchild 1.0 to be %v, but got %v", expectedPath, actualPath)
	}
}

func TestTree_MoveCursorDown(t *testing.T) {
	tree := newTestTree()
	testCases := []struct {
		name     string
		start    []int
		expected []int
	}{
		{"from empty, select first", []int{}, []int{0}},
		{"from first, enter expanded", []int{0}, []int{0, 0}},
		{"from grandchild, to sibling", []int{0, 0}, []int{0, 1}},
		{"from last grandchild, to next uncle", []int{0, 1}, []int{1}},
		{"from collapsed, to sibling", []int{1}, []int{2}},
		{"from last node, stay", []int{2}, []int{2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree.SelectionPath = tc.start
			tree.MoveCursorDown()
			if !reflect.DeepEqual(tree.SelectionPath, tc.expected) {
				t.Errorf("Expected path %v, but got %v", tc.expected, tree.SelectionPath)
			}
		})
	}
}

func TestTree_MoveCursorUp(t *testing.T) {
	tree := newTestTree()
	testCases := []struct {
		name     string
		start    []int
		expected []int
	}{
		// Visible tree structure from newTestTree():
		// - root (expanded)		(Y=N/A)
		//   - child 0 (expanded)	(Y=0)
		//   - grandchild 0.0		(Y=1)
		//   - grandchild 0.1		(Y=2)
		//   - child 1 (collapsed)	(Y=3)
		//   - grandchild 1.0		(Y=N/A)
		//   - child 2				(Y=4)
		{"from empty, stay", []int{}, []int{}},
		{"from first, to parent (root)", []int{0}, []int{}},
		{"from child 1, to sibling (child 0) (expanded to grandchild 0.1)", []int{1}, []int{0, 1}},
		{"from grandchild 0.1, to sibling grandchild 0.0", []int{0, 1}, []int{0, 0}},
		{"from first grandchild 0.0, to parent (child 0)", []int{0, 0}, []int{0}},
		{"from child 2, to sibling (child 1)", []int{2}, []int{1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree.SelectionPath = tc.start
			tree.MoveCursorUp()
			if !reflect.DeepEqual(tree.SelectionPath, tc.expected) {
				t.Errorf("Expected path %v, but got %v", tc.expected, tree.SelectionPath)
			}
		})
	}
}

func TestTree_MoveCursorLeft(t *testing.T) {
	// tree := newTestTree()
	testCases := []struct {
		name           string
		startPath      []int
		expectedPath   []int
		expectedIsNode bool // true if we expect a node, false if we expect it to be collapsed
		expectedState  bool // expanded or collapsed state
	}{
		{"collapse expanded node", []int{0}, []int{0}, true, false},             // Stays on the node but collapses it
		{"move from collapsed to parent", []int{1}, []int{}, false, false},      // Moves to parent (root)
		{"move from leaf to parent", []int{0, 0}, []int{0}, false, false},       // Moves to parent
		{"move from collapsed leaf to parent", []int{2}, []int{}, false, false}, // Moves to parent (root)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset tree state for each test
			testTree := newTestTree()
			testTree.SelectionPath = tc.startPath
			testTree.MoveCursorLeft()

			if !reflect.DeepEqual(testTree.SelectionPath, tc.expectedPath) {
				t.Errorf("Expected path %v, got %v", tc.expectedPath, testTree.SelectionPath)
			}

			if tc.expectedIsNode {
				node := testTree.GetAtPath(tc.expectedPath)
				if node.Expanded != tc.expectedState {
					t.Errorf("Expected node at %v to have Expanded=%t, but was %t", tc.expectedPath, tc.expectedState, node.Expanded)
				}
			}
		})
	}
}

func TestTree_MoveCursorRight(t *testing.T) {
	// tree := newTestTree()
	testCases := []struct {
		name          string
		startPath     []int
		expectedPath  []int
		expectedState bool // expanded state
	}{
		{"expand collapsed node", []int{1}, []int{1}, true},
		{"do nothing on expanded node", []int{0}, []int{0}, true},
		{"do nothing on leaf node", []int{2}, []int{2}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset tree state for each test
			testTree := newTestTree()
			testTree.SelectionPath = tc.startPath
			testTree.MoveCursorRight()

			if !reflect.DeepEqual(testTree.SelectionPath, tc.expectedPath) {
				t.Errorf("Expected path %v, got %v", tc.expectedPath, testTree.SelectionPath)
			}

			node := testTree.GetAtPath(tc.expectedPath)
			if node.Expanded != tc.expectedState {
				t.Errorf("Expected node at %v to have Expanded=%t, but was %t", tc.expectedPath, tc.expectedState, node.Expanded)
			}
		})
	}
}

func TestTree_ToggleExpand(t *testing.T) {
	tree := newTestTree()

	// Start on node {1}, which is collapsed
	tree.SelectionPath = []int{1}
	node := tree.GetAtPath(tree.SelectionPath)
	if node.Expanded != false {
		t.Fatalf("Prerequisite failed: node {1} should be collapsed.")
	}

	// Toggle 1: should expand
	tree.ToggleExpand()
	if node.Expanded != true {
		t.Errorf("Expected node to be expanded after first toggle, but it was not.")
	}

	// Toggle 2: should collapse
	tree.ToggleExpand()
	if node.Expanded != false {
		t.Errorf("Expected node to be collapsed after second toggle, but it was not.")
	}

	// Test on a leaf node (should do nothing)
	tree.SelectionPath = []int{2}
	leafNode := tree.GetAtPath(tree.SelectionPath)
	leafNode.Expanded = false // ensure state
	tree.ToggleExpand()
	if leafNode.Expanded != false {
		t.Errorf("Toggling a leaf node should not change its expanded state.")
	}
}

func TestTree_calculateCursorY(t *testing.T) {
	tree := newTestTree()
	testCases := []struct {
		name     string
		path     []int
		expected int
	}{
		// Visible tree structure from newTestTree():
		// - root (expanded)		(Y=N/A)
		//   - child 0 (expanded)	(Y=0)
		//   - grandchild 0.0		(Y=1)
		//   - grandchild 0.1		(Y=2)
		//   - child 1 (collapsed)	(Y=3)
		//   - grandchild 1.0		(Y=N/A)
		//   - child 2				(Y=4)
		{"first child (child 0)", []int{0}, 0},
		{"first grandchild (0.0)", []int{0, 0}, 1},
		{"second grandchild (0.1)", []int{0, 1}, 2},
		{"collapsed node (child 1)", []int{1}, 3},
		{"last child (child 2)", []int{2}, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree.SelectionPath = tc.path
			y := tree.CursorY()
			if y != tc.expected {
				t.Errorf("Expected Y position %d for path %v, but got %d", tc.expected, tc.path, y)
			}
		})
	}
}
