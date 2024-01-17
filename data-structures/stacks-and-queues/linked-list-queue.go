package main

import "fmt"

type Queue struct {
	First  *Node
	Last   *Node
	Length int
}

func (q *Queue) Peek() *Node {
	return q.First
}

func (q *Queue) Enqueue(value string) {
	newNode := &Node{
		Value: value,
		Next:  q.Last,
	}
	q.Last = newNode
	if q.Length == 0 {
		q.First = newNode
	}
	q.Length++
}

func (q *Queue) Dequeue() {
	if q.First == q.Last {
		q.First = nil
		q.Last = nil
		return
	}

	newFirst := q.Last
	for ; newFirst.Next != q.First; newFirst = newFirst.Next {
	}
	newFirst.Next = nil
	q.First = newFirst
}

func (q *Queue) String() string {
	var str string
	for node := q.Last; node != nil; node = node.Next {
		str += node.Value
		if node.Next != nil {
			str += ", "
		}
	}
	return fmt.Sprintf("[%s]", str)
}

func testLinkedListQueue() {
	var llq Queue

	fmt.Println("LinkedList Queue")
	fmt.Println(llq.String())
	llq.Enqueue("Foo")
	fmt.Println(llq.String())
	llq.Enqueue("Bar")
	llq.Enqueue("Baz")
	fmt.Println(llq.String())
	llq.Dequeue()
	fmt.Println(llq.String())
	llq.Dequeue()
	llq.Dequeue()
	fmt.Println(llq.String())
}
