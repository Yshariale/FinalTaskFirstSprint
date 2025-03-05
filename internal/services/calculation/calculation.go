package calculation

// 
// 
// В этом модуле расположена логика работы с БИНАРНЫМ ДЕРЕВОМ и ОБРАТНОЙ ПОЛЬСКОЙ НОТАЦИЕЙ
// 
// 

import (
	"strconv"
	"strings"
	"unicode"
)

type Tree struct {
	Root *TreeNode
}

type TreeNode struct {
	Val    string
	Left   *TreeNode
	Right  *TreeNode
	TaskID int
}

// Проверка на готовность функции родить задачу. 
// Если у вершины оба потомка - числа, то вершина готова
func (node *TreeNode) IsSpare() bool {
	if node.Right != nil && node.Left != nil {
		_, err1 := strconv.ParseFloat(node.Left.Val, 64)
		_, err2 := strconv.ParseFloat(node.Right.Val, 64)
		if err1 == nil && err2 == nil {
			return true
		}
	}
	return false
}

// Поиск всех вершин, готовых родить задачу
func (t *Tree) FindSpareNodes() []*TreeNode {
	spare_nodes := []*TreeNode{}
	stack := []*TreeNode{t.Root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.IsSpare() {
			spare_nodes = append(spare_nodes, node)
		} else {
			if node.Right != nil {
				stack = append(stack, node.Right)
			}
			if node.Left != nil {
				stack = append(stack, node.Left)
			}
		}
	}
	return spare_nodes
}

// Когда задача решена, заменяем вершину на просто число,
// чтобы ее родительская вершина была готова родить задачу
func (t *Tree) ReplaceNodeWithValue(node *TreeNode, val float64) {
	node.Left = nil
	node.Right = nil
	arg := strconv.FormatFloat(val, 'f', -1, 64)
	node.Val = arg
}

// Поиск родительской вершины и вершины по ID задачи. 
// Нужно для замены вершины на число после решения задачи
func (t *Tree) FindParentAndNodeByTaskID(task_id int) (*TreeNode, *TreeNode) {
	if t.Root.TaskID == task_id {
		return nil, t.Root
	}

	stack := []*TreeNode{t.Root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.Right != nil && node.Right.TaskID == task_id {
			return node, node.Right
		}
		if node.Left != nil && node.Left.TaskID == task_id {
			return node, node.Left
		}

		if node.Right != nil {
			stack = append(stack, node.Right)
		}
		if node.Left != nil {
			stack = append(stack, node.Left)
		}

	}
	return nil, nil
}

// Переводит из инфиксной в постфиксную запись (знаю умные слова)
// А еще по пути проверяет выражение на валидность
func ToPostfix(expression string) ([]string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	var output []string
	var stack []string

	priority := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	if len(expression) == 0 {
		return nil, ErrInvalidExpression
	}

	var prevToken string

	for i := 0; i < len(expression); i++ {
		char := string(expression[i])

		if unicode.IsDigit(rune(expression[i])) || char == "." ||
			(char == "-" && (i == 0 || prevToken == "(" || priority[prevToken] > 0)) {

			number := char
			for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || string(expression[i+1]) == ".") {
				i++
				number += string(expression[i])
			}
			output = append(output, number)
			prevToken = number

		} else if char == "(" {
			stack = append(stack, char)
			prevToken = char

		} else if char == ")" {
			if prevToken == "" || priority[prevToken] > 0 || prevToken == "(" {
				return nil, ErrMismatchedBracket
			}

			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			if len(stack) == 0 {
				return nil, ErrMismatchedBracket
			}

			stack = stack[:len(stack)-1]
			prevToken = ")"

		} else if priority[char] > 0 {
			if prevToken == "" || priority[prevToken] > 0 || prevToken == "(" {
				return nil, ErrInvalidOperationsPlacement
			}

			for len(stack) > 0 && priority[stack[len(stack)-1]] >= priority[char] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
			prevToken = char

		} else {
			return nil, ErrInvalidSymbols
		}
	}

	if prevToken == "" || priority[prevToken] > 0 {
		return nil, ErrInvalidExpression
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, ErrMismatchedBracket
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// Строит бинарное дерево из постфиксной записи
func BuildTree(postfix []string) *Tree {
	stack := []*TreeNode{}

	for _, token := range postfix {
		_, err := strconv.ParseFloat(token, 64)
		if err == nil {
			stack = append(stack, &TreeNode{Val: token})
		} else {
			if len(stack) < 2 {
				panic("Invalid expression: not enough operands")
			}
			right := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			node := &TreeNode{Val: token, Left: left, Right: right}
			stack = append(stack, node)
		}
	}

	if len(stack) != 1 {
		panic("Invalid expression: leftover nodes in stack")
	}

	return &Tree{Root: stack[0]}
}
