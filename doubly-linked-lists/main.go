package main

import "fmt"

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head   *Node
	tail   *Node
	length int
}

func New(value int) *LinkedList {
	head := Node{
		value: value,
		next:  nil,
	}

	return &LinkedList{
		head:   &head,
		tail:   &head,
		length: 1,
	}
}

func (ll *LinkedList) append(value int) {
	newNode := Node{
		value: value,
		next:  nil,
	}

	ll.tail.next = &newNode
	ll.tail = &newNode
	ll.length++
}

func (ll *LinkedList) prepend(value int) {
	newNode := Node{
		value: value,
		next:  ll.head,
	}

	ll.head = &newNode
	ll.length += 1
}

func (ll *LinkedList) String() string {
	var str string = "["
	for node := ll.head; node != nil; node = node.next {
		str += fmt.Sprint(node.value)

		if node.next != nil {
			str += ", "
		}
	}
	str += "]"
	return str
}

func (ll *LinkedList) insert(index int, value int) error {
	if index > ll.length {
		ll.append(value)
		return nil
	} else if index < 0 {
		return fmt.Errorf("IndexOutOfBounds")
	}
	prevNode := ll.traverseToIndex(index - 1)
	nextNode := prevNode.next
	newNode := Node{
		value: value,
		next:  nextNode,
	}
	prevNode.next = &newNode
	ll.length++
	return nil
}

func (ll *LinkedList) traverseToIndex(index int) *Node {
	i := 0
	node := ll.head
	for ; i != index; node = node.next {
		i++
	}
	return node
}

func (ll *LinkedList) remove(index int) error {
	if index > ll.length {
		return fmt.Errorf("No")
	}

	prevNode := ll.traverseToIndex(index - 1)
	unwantedNode := prevNode.next
	prevNode.next = unwantedNode.next
	ll.length--
	return nil
}

func main() {
	list := New(10)
	fmt.Println(list.String())
	list.append(5)
	fmt.Println(list.String())
	list.append(16)
	fmt.Println(list.String())
	list.prepend(1)
	fmt.Println(list.String())
	list.insert(2, 99)
	fmt.Println(list.String())
	list.insert(20, 88)
	fmt.Println(list.String())
	list.remove(2)
	fmt.Println(list.String())
}
