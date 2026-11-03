package main

import (
	"errors"
	"os"
	"reflect"
)

type Node struct {
	prev  *Node
	next  *Node
	value int
}

type LinkedList2 struct {
	head *Node
	tail *Node
}

func (l *LinkedList2) AddInTail(item Node) {
	if l.head == nil {
		l.head = &item
    l.head.next = nil
		l.head.prev = nil
	} else {
		l.tail.next = &item
    item.prev = l.tail
	}

	l.tail = &item
  l.tail.next = nil
}

func (l *LinkedList2) Count() int {
	count := 0
	node := l.head
	for node != nil {
		count++
		node = node.next
	}

	return count
}

func (l *LinkedList2) Find(n int) (Node, error) {
	node := l.head
	for node != nil {
		if node.value == n {
			return *node, nil
		}
		node = node.next
	}

	return Node{}, errors.New("not found")
}

func (l *LinkedList2) FindAll(n int) []Node {
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

func (l *LinkedList2) Delete(n int, all bool) {
	// Empty list
	if l.head == nil {
		return
	}

	node := l.head
	for node != nil {
		if node.value == n {
			// Single-element list
			if node.next == nil && node.prev == nil {
				l.head = nil
				l.tail = nil
	
				return
			}

			switch node {
			case l.head:
				node.next.prev = nil
				l.head = node.next
			case l.tail:
				node.prev.next = nil
				l.tail = node.prev

				return
			default:
				prev := node.prev
				next := node.next
				prev.next = next
				next.prev = prev
			}

			node = node.next

			if !all {
				return
			}
		} else {
			node = node.next
		}
	}
}

func (l *LinkedList2) Insert(after *Node, add Node) {
	oldNext := after.next
	
	add.next = after.next
	add.prev = after
	after.next = &add

	if after == l.tail {
		l.tail = &add
		return
	}

	oldNext.prev = &add
}

func (l *LinkedList2) InsertFirst(first Node) {
	if l.head == nil {
		l.head = &first
		l.tail = &first

		return
	}

	l.head.prev = &first
	first.next = l.head
	l.head = &first
}

func (l *LinkedList2) Clean() {
	l.head = nil
	l.tail = nil
}


