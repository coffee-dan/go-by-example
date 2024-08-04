package main

import "fmt"

/**
 * Definition for singly-linked list.
 */
type ListNode struct {
	Val  int
	Next *ListNode
}

func (l *ListNode) String() (output string) {
	for node := l; node != nil; node = node.Next {
		output += fmt.Sprintf("%d->", node.Val)
	}
	return output
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var head, prev *ListNode
	head = nil
	prev = head
	node1 := l1
	node2 := l2
	var val, remainder, quotient int
	for {
		val = quotient

		if val == 0 && node1 == nil && node2 == nil {
			break
		}

		if node1 != nil {
			val += node1.Val
			fmt.Printf("%d + ", node1.Val)
			node1 = node1.Next
		}

		if node2 != nil {
			val += node2.Val
			fmt.Printf("%d = ", node2.Val)
			node2 = node2.Next
		}

		quotient = val / 10
		remainder = val % 10

		curr := &ListNode{remainder, nil}
		if prev == nil {
			head = curr
		} else {
			prev.Next = curr
		}
		prev = curr
	}

	return head
}

func intToListNode(val int64) (head *ListNode) {
	head = &ListNode{}
	head.Val = int(val % 10)
	val = val / 10

	prev := head
	for n := val; n != 0; {
		rem := int(n % 10)
		n = n / 10

		curr := &ListNode{rem, nil}
		prev.Next = curr
		prev = curr
	}

	return head
}

func main() {
	l1 := intToListNode(int64(456))
	fmt.Println(l1)
	l2 := intToListNode(int64(465))
	fmt.Println(l2)
	fmt.Println(addTwoNumbers(l1, l2))
}
