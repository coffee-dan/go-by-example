package main

import "fmt"

type Stack struct {
	Top    *Node
	Bottom *Node
	Length int
}

func (s *Stack) Peek() *Node {
	return s.Top
}

func (s *Stack) Push(value string) {
	newNode := &Node{
		Value: value,
		Next:  s.Top,
	}
	s.Top = newNode

	if s.Length == 0 {
		s.Bottom = newNode
	}

	s.Length++
}

func (s *Stack) Pop() *Node {
	if s.Top == nil {
		return nil
	}

	oldTop := s.Top
	s.Top = s.Top.Next

	if s.Top == s.Bottom {
		s.Bottom = nil
	}

	s.Length--
	return oldTop
}

func (s *Stack) String() string {
	str := "["
	for node := s.Top; node != nil; node = node.Next {
		str += fmt.Sprintf("\"%s\"", node.Value)
		if node.Next != nil {
			str += ", "
		}
	}
	str += "]"
	return str
}

func testLinkedListStack() {
	var lls Stack
	fmt.Println("LinkedList Stack")
	lls.Pop()
	lls.Push("Google")
	lls.Push("Udemy")
	lls.Push("Discord")
	fmt.Println(lls.String())
	fmt.Println(lls.Peek())
	lls.Pop()
	fmt.Println(lls.String())
	lls.Pop()
	lls.Pop()
	fmt.Println(lls.String())
}
