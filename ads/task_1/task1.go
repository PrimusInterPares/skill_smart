package main

import(
	"errors"
//	"os"
//	"reflect"
)

type Node struct {
	next  *Node
	value int
}

type LinkedList struct {
	head *Node
	tail *Node
}

func (l *LinkedList) AddInTail(item Node) {
	if l.head == nil {
		l.head = &item
	} else {
		l.tail.next = &item
	}

	l.tail = &item
}

func (l *LinkedList) Count() int {
	count := 0
	node := l.head
	for node != nil {
		count++
		node = node.next
	}

	return count
}

func (l *LinkedList) Find(n int) (Node, error) {
	node := l.head
	for node != nil {
		if node.value == n {
			return *node, nil
		}
		node = node.next
	}

	return Node{}, errors.New("not found")
}

func (l *LinkedList) FindAll(n int) []Node {
	var nodes []Node

	node := l.head
	for node != nil {
		if node.value == n {
			nodes = append(nodes, *node)
		}

		node = node.next
	}

	return nodes
}

func (l *LinkedList) Delete(n int, all bool) {
	// Empty list
	if l.head == nil {
		return
	}

	// Delete sequentially matching heads
	for l.head != nil && l.head.value == n {
		l.head = l.head.next

		// List is empty
		if l.head == nil {
			l.tail = nil
			return
		}
		
		if !all {
			// One element left
			if l.head != nil && l.head.next == nil {
				l.tail = l.head
			}
			return
		}
	}

	currNode := l.head
	for currNode.next != nil {
		if currNode.next.value == n {
			if currNode.next == l.tail {
				l.tail = currNode
				currNode.next = nil
				return // End of the list
			}

			currNode.next = currNode.next.next
			
			if !all {
				return
			}
		} else {
			currNode = currNode.next
		}
	}
}

//func (l *LinkedList) Insert(after *Node, add Node) {
//}

//func (l *LinkedList) InsertFirst(first Node) {
//}

func (l *LinkedList) Clean() {
	l.head = nil
	l.tail = nil
}


