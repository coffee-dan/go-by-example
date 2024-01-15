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

func NewLinkedList(value int) *LinkedList {
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

type DoublyLinkedNode struct {
	value int
	next  *DoublyLinkedNode
	prev  *DoublyLinkedNode
}

type DoublyLinkedList struct {
	head   *DoublyLinkedNode
	tail   *DoublyLinkedNode
	length int
}

func NewDoublyLinkedList(value int) *DoublyLinkedList {
	head := DoublyLinkedNode{
		value: value,
		next:  nil,
	}

	return &DoublyLinkedList{
		head:   &head,
		tail:   &head,
		length: 1,
	}
}

func (dll *DoublyLinkedList) append(value int) {
	newNode := DoublyLinkedNode{
		value: value,
		next:  nil,
		prev:  dll.tail,
	}

	dll.tail.next = &newNode
	dll.tail = &newNode
	dll.length++
}

func (dll *DoublyLinkedList) prepend(value int) {
	newNode := DoublyLinkedNode{
		value: value,
		next:  dll.head,
		prev:  nil,
	}

	dll.head.prev = &newNode
	dll.head = &newNode
	dll.length++
}

func (dll *DoublyLinkedList) String() string {
	var str string = "["
	for node := dll.head; node != nil; node = node.next {
		str += fmt.Sprint(node.value)

		if node.next != nil {
			str += ", "
		}
	}
	str += "]"
	return str
}

func (dll *DoublyLinkedList) insert(index int, value int) error {
	if index > dll.length {
		dll.append(value)
		return nil
	} else if index < 0 {
		return fmt.Errorf("IndexOutOfBounds")
	}
	prevNode := dll.traverseToIndex(index - 1)
	nextNode := prevNode.next
	newNode := DoublyLinkedNode{
		value: value,
		next:  nextNode,
		prev:  prevNode,
	}
	prevNode.next = &newNode
	nextNode.prev = &newNode
	dll.length++
	return nil
}

func (dll *DoublyLinkedList) traverseToIndex(index int) *DoublyLinkedNode {
	i := 0
	node := dll.head
	for ; i != index; node = node.next {
		i++
	}
	return node
}

func (dll *DoublyLinkedList) remove(index int) error {
	if index > dll.length {
		return fmt.Errorf("No")
	}

	prevNode := dll.traverseToIndex(index - 1)
	unwantedNode := prevNode.next
	nextNode := unwantedNode.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
	dll.length--
	return nil
}

func testLinkedList() string {
	list := NewLinkedList(10)
	list.append(5)
	list.append(16)
	list.prepend(1)
	list.insert(2, 99)
	list.insert(20, 88)
	list.remove(2)
	return list.String()
}

func testDoublyLinkedList() string {
	list := NewDoublyLinkedList(10)
	list.append(5)
	list.append(16)
	list.prepend(1)
	list.insert(2, 99)
	list.insert(20, 88)
	list.remove(2)
	return list.String()
}

func (ll *LinkedList) CloneReversed() *LinkedList {
	revLL := LinkedList{}

	for node := ll.head; node != nil; node = node.next {
		revLL.prepend(node.value)
	}

	return &revLL
}

func (ll *LinkedList) Reverse() {
	prev := ll.head
	ll.head = ll.tail
	ll.tail = prev

	curr := prev.next
	for curr != nil {
		next := curr.next
		curr.next = prev
		prev = curr
		curr = next
	}

	ll.tail.next = nil
}

func main() {
	fmt.Println(testLinkedList())
	fmt.Println(testDoublyLinkedList())

	ll := NewLinkedList(1)
	ll.append(2)
	ll.append(3)
	ll.append(4)

	fmt.Println(ll.String())

	rll := ll.CloneReversed()

	fmt.Println(rll.String())
	fmt.Println(rll.length)

	ll.Reverse()

	fmt.Println(ll.String())
	fmt.Println(ll.length)
}
