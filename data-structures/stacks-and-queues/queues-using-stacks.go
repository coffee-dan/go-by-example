package main

type IntNode struct {
	Value int
	Next  *IntNode
}

type IntStack struct {
	Top    *IntNode
	Length int
}

type MyQueue struct {
	in  *IntStack
	out *IntStack
}

func Constructor() MyQueue {
	return MyQueue{
		in:  &IntStack{nil, 0},
		out: &IntStack{nil, 0},
	}
}

func (sta *IntStack) Push(x int) {
	sta.Top = &IntNode{
		Value: x,
		Next:  sta.Top,
	}
	sta.Length++
}

func (sta *IntStack) Pop() (val int) {
	val = sta.Top.Value
	sta.Top = sta.Top.Next
	sta.Length--
	return
}

func (sta *IntStack) Empty() bool {
	return sta.Length == 0
}

func (que *MyQueue) Push(x int) {
	que.in.Push(x)
}

func (que *MyQueue) Pop() int {
	if que.out.Empty() {
		for !que.in.Empty() {
			que.out.Push(que.in.Pop())
		}
	}

	return que.out.Pop()
}

func (que *MyQueue) Peek() int {
	if que.out.Empty() {
		for !que.in.Empty() {
			que.out.Push(que.in.Pop())
		}
	}
	return que.out.Top.Value
}

func (que *MyQueue) Empty() bool {
	return que.in.Empty() && que.out.Empty()
}
