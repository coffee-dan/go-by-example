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
		if node.Value > value {
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

func (bst *BinarySearchTree) Remove(value int) *TreeNode {
	node := bst.Root
	var junk *TreeNode
	var replacement *TreeNode
	for {
		if node.Value == value {
			junk = node
			replacement = node.ReplaceWithSuccessor()
			break
		} else if node.Right != nil && node.Right.Value == value {
			junk = node.Right
			replacement = junk.ReplaceWithSuccessor()
			node.Right = replacement
			break
		} else if node.Left != nil && node.Left.Value == value {
			junk = node.Left
			replacement = junk.ReplaceWithSuccessor()
			node.Left = replacement
			break
		} else if node.Value > value {
			node = node.Left
		} else if node.Value < value {
			node = node.Right
		} else {
			break
		}
	}

	if junk == bst.Root {
		bst.Root = replacement
	}
	return junk
}

func (tr *TreeNode) ReplaceWithSuccessor() *TreeNode {
	if tr.Left == nil && tr.Right == nil {
		return nil
	}
	parent := tr
	var replacement *TreeNode
	if tr.Right != nil {
		replacement = tr.Right
		for {
			if replacement.Left == nil && replacement.Right == nil {
				if parent == tr {
					parent.Right = nil
				} else {
					parent.Left = nil
				}

				replacement.Left = tr.Left
				replacement.Right = tr.Right
				break
			}
			parent = replacement
			replacement = replacement.Right
		}
	} else if tr.Left != nil {
		replacement = tr.Left
		for {
			if replacement.Left == nil && replacement.Right == nil {
				if parent == tr {
					parent.Left = nil
				} else {
					parent.Right = nil
				}

				replacement.Left = tr.Left
				replacement.Right = tr.Right
				break
			}
			parent = replacement
			replacement = replacement.Left
		}
	}

	return replacement
}

func (tr *TreeNode) renderSubtree(spacer string) string {
	var rightStr string
	var leftStr string
	if tr.Right != nil {
		rightStr = spacer + "├── " + tr.Right.renderSubtree(spacer+"│   ")
	} else {
		rightStr = spacer + "├──*"
	}

	if tr.Left != nil {
		leftStr = spacer + "└── " + tr.Left.renderSubtree(spacer+"    ")
	} else {
		leftStr = spacer + "└──*"
	}

	str := fmt.Sprintf("%02d\n%s\n%s", tr.Value, rightStr, leftStr)
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
	fmt.Println(bst.String())

	bst.Remove(20)
	fmt.Println(bst.String())
	bst.Remove(15)
	fmt.Println(bst.String())
	bst.Remove(9)
	fmt.Println(bst.String())
	bst.Remove(4)
	fmt.Println(bst.String())
}
