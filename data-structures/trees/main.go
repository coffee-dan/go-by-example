package main

import "fmt"

type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

type BinarySearchTree struct {
	Root *TreeNode
}

func (bst *BinarySearchTree) Insert(value int) {
	newNode := &TreeNode{value, nil, nil}

	if bst.Root == nil {
		bst.Root = newNode
		return
	}

	node := bst.Root
	for {
		if node.Value < value {
			if node.Left == nil {
				node.Left = newNode
				return
			} else {
				node = node.Left
			}
		} else {
			if node.Right == nil {
				node.Right = newNode
				return
			} else {
				node = node.Right
			}
		}
	}
}

func (bst *BinarySearchTree) Lookup(value int) *TreeNode {
	node := bst.Root
	for {
		if node.Value < value {
			node = node.Left
		} else if node.Value > value {
			node = node.Right
		} else {
			return node
		}
	}
}

func (tr *TreeNode) renderSubtree(spacer string) string {
	var leftStr string
	var rightStr string
	if tr.Left != nil {
		leftStr = spacer + "├── " + tr.Left.renderSubtree(spacer+"│   ")
	} else {
		leftStr = spacer + "├──*"
	}

	if tr.Right != nil {
		rightStr = spacer + "└── " + tr.Right.renderSubtree(spacer+"    ")
	} else {
		rightStr = spacer + "└──*"
	}

	str := fmt.Sprintf("%02d\n%s\n%s", tr.Value, leftStr, rightStr)
	return str
}

func (bst *BinarySearchTree) String() string {
	return bst.Root.renderSubtree("")
}

func main() {
	bst := BinarySearchTree{}
	bst.Insert(9)
	bst.Insert(4)
	bst.Insert(6)
	bst.Insert(20)
	bst.Insert(170)
	bst.Insert(15)
	bst.Insert(1)
	// bst.Insert(1)
	fmt.Println(bst.String())

	fmt.Println(bst.Lookup(170).renderSubtree(""))
	fmt.Println(bst.Lookup(20).renderSubtree(""))
}
