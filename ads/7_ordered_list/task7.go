package main

import (
	"constraints"
	"errors"
	// "os"
	"strings"
)

var (
	ErrEmptyList    = errors.New("list is empty")
	ErrNodeNotFound = errors.New("node is not found")
)

type Node[T constraints.Ordered] struct {
	prev  *Node[T]
	next  *Node[T]
	value T
}

type OrderedList[T constraints.Ordered] struct {
	head      *Node[T]
	tail      *Node[T]
	ascending bool
	count     int
}

func (l *OrderedList[T]) Count() int {
	return l.count
}

func (l *OrderedList[T]) Add(item T) {
	newNode := &Node[T]{value: item}

	if l.head == nil {
		l.head = newNode
		l.tail = newNode
		l.count++

		return
	}

	node := l.head
	for node != nil {
		if l.Compare(item, node.value) < 0 {
			newNode.next = node
			newNode.prev = node.prev

			if node.prev != nil {
				node.prev.next = newNode
			} else {
				l.head = newNode
			}
			node.prev = newNode

			l.count++

			return
		}

		node = node.next
	}

	newNode.prev = l.tail
	l.tail.next = newNode
	l.tail = newNode

	l.count++
}

func (l *OrderedList[T]) Find(n T) (Node[T], error) {
	var zero Node[T]

	if l.Count() == 0 {
		return zero, ErrEmptyList
	}

	node := l.head
	for node != nil {
		compare := l.Compare(node.value, n)

		if compare > 0 {
			break
		}

		if compare == 0 {
			return *node, nil
		}

		node = node.next
	}

	return zero, ErrNodeNotFound
}

func (l *OrderedList[T]) Delete(n T) {
	node := l.head
	for node != nil {
		if node.value == n {
			// single-element list
			if node.next == nil && node.prev == nil {
				l.head = nil
				l.tail = nil
				l.count--

				return
			}

			switch node {
			case l.head:
				node.next.prev = nil
				l.head = node.next
			case l.tail:
				node.prev.next = nil
				l.tail = node.prev
			default:
				prev := node.prev
				next := node.next
				prev.next = next
				next.prev = prev
			}

			l.count--

			return
		}

		node = node.next
	}
}

func (l *OrderedList[T]) Clear(asc bool) {
	l.head = nil
	l.tail = nil
	l.count = 0
	l.ascending = asc
}

func (l *OrderedList[T]) Compare(v1, v2 T) int {
	result := 0

	if s1, ok := any(v1).(string); ok {
		s2 := any(v2).(string)

		s1, s2 = strings.TrimSpace(s1), strings.TrimSpace(s2)

		if s1 < s2 {
			result = -1
		} else if s1 > s2 {
			result = +1
		}

		if !l.ascending {
			return -result
		}

		return result
	}

	if v1 < v2 {
		result = -1
	} else if v1 > v2 {
		result = +1
	}

	if !l.ascending {
		return -result
	}

	return result
}
