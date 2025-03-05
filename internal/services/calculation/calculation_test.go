package calculation

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsSpare(t *testing.T) {
	node := &TreeNode{Val: "+", Left: &TreeNode{Val: "3"}, Right: &TreeNode{Val: "4"}}
	assert.True(t, node.IsSpare(), "Expected node to be spare")

	node = &TreeNode{Val: "*", Left: &TreeNode{Val: "*"}, Right: &TreeNode{Val: "5"}}
	assert.False(t, node.IsSpare(), "Expected node not to be spare")
}

func TestFindSpareNodes(t *testing.T) {
	tree := &Tree{
		Root: &TreeNode{
			Val: "+",
			Left:  &TreeNode{Val: "3"},
			Right: &TreeNode{Val: "4"},
		},
	}

	spareNodes := tree.FindSpareNodes()
	assert.Len(t, spareNodes, 1, "Expected 1 spare node")
}

func TestReplaceNodeWithValue(t *testing.T) {
	tree := &Tree{Root: &TreeNode{Val: "+", Left: &TreeNode{Val: "3"}, Right: &TreeNode{Val: "4"}}}
	tree.ReplaceNodeWithValue(tree.Root, 7)
	if tree.Root.Val != "7" || tree.Root.Left != nil || tree.Root.Right != nil {
		t.Errorf("Node replacement failed")
	}
}

func TestFindParentAndNodeByTaskID(t *testing.T) {
	tree := &Tree{
		Root: &TreeNode{
			TaskID: 1,
			Left:  &TreeNode{TaskID: 2},
			Right: &TreeNode{TaskID: 3},
		},
	}

	parent, node := tree.FindParentAndNodeByTaskID(2)
	if node == nil || node.TaskID != 2 || parent.TaskID != 1 {
		t.Errorf("Parent or node lookup failed")
	}
}

func TestToPostfix(t *testing.T) {
	t.Run("Valid expressions", func(t *testing.T) {
		for _, test := range ValidTestSet {
			t.Run(test.Name, func(t *testing.T) {
				result, err := ToPostfix(test.Expression)
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, test.Expected_answer, result, "Incorrect postfix notation")
			})
		}
	})
	t.Run("Invalid expressions", func(t *testing.T) {
		for _, test := range InvalidTestSet {
			t.Run(test.Name, func(t *testing.T) {
				_, err := ToPostfix(test.Expression)
				assert.Error(t, err, "Expected error")
				assert.ErrorIs(t, test.Expected_error, err, "Incorrect error")
			})
		}
	})
}

func TestBuildTree(t *testing.T) {
	postfix := []string{"3", "4", "+"}
	tree := BuildTree(postfix)
	if tree.Root.Val != "+" || tree.Root.Left.Val != "3" || tree.Root.Right.Val != "4" {
		t.Errorf("Tree building failed")
	}
}
