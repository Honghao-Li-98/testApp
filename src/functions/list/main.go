package main

import (
	"fmt"
)

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	head *Node[T]
}

type Node[T any] struct {
	val  T
	next *Node[T]
}

func (list *List[T]) Push(v T) {
	if list.head == nil {
		list.head = &Node[T]{val: v}
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = &Node[T]{val: v}
	}
}

func (list *List[T]) RemoveLastItem() {
	if list.head == nil {
		return
	} else if list.head.next == nil {
		list.head = nil
	} else {
		current := list.head

		for current.next.next != nil {
			current = current.next
		}

		current.next = nil
	}
}

func (list *List[T]) Print() {

	current := list.head

	if current == nil {
		fmt.Println("Empty List")
		return
	}

	for current != nil {
		fmt.Println(current.val)
		current = current.next
	}
}

func main() {
	list := List[int]{}

	list.Push(1)
	list.Push(2)
	list.Push(3)
	list.Push(4)
	list.Push(5)

	list.RemoveLastItem()

	list.Print()
}
