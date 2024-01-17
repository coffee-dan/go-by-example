package main

import "fmt"

func Peek(stack []string) string {
	return stack[len(stack)-1]
}

func Push(stack []string, value string) []string {
	return append(stack, value)
}

func Pop(stack []string) []string {
	if len(stack) > 0 {
		return stack[:len(stack)-1]
	}
	return stack
}

func testArrayStack() {
	var ars []string
	fmt.Println("Array Stack")
	ars = Pop(ars)
	ars = Push(ars, "Google")
	ars = Push(ars, "Udemy")
	ars = Push(ars, "Discord")
	fmt.Println(ars)
	fmt.Println(Peek(ars))
	ars = Pop(ars)
	fmt.Println(ars)
	ars = Pop(ars)
	ars = Pop(ars)
	fmt.Println(ars)
}
