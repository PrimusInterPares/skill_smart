package main

import (
	"constraints"
	"errors"
)

// Порядковый номер задания на курсе: 7
// Номер задачи из задания: 8
// Краткое название: метод удаления всех дубликатов из упорядоченного списка
// Time complexity: O(N)
// Memory complexity: O(1)
var ErrOrderingMismatch = errors.New("ordering must match for both lists")

func (l *OrderedList[T]) RemoveDuplicates() {
	if l.Count() < 2 {
		return
	}

	node := l.head
	next := node.next
	for next != nil {
		if l.Compare(next.value, node.value) == 0 {
			node.next = next.next

			if next.next != nil {
				next.next.prev = node
			} else {
				l.tail = node
			}

			l.count--

			next = node.next
		} else {
			node = node.next
			next = node.next
		}
	}
}

// Номер задачи из задания: 9
// Краткое название: алгоритм слияния двух упорядоченных списков в один, сохраняя порядок элементов
// Time complexity: O(N)
// Memory complexity: O(N)
func MergeOrderedLists[T constraints.Ordered](list1, list2 *OrderedList[T]) (*OrderedList[T], error) {
	if list1.ascending != list2.ascending {
		return &OrderedList[T]{}, ErrOrderingMismatch
	}

	merged := &OrderedList[T]{ascending: list1.ascending}

	node1 := list1.head
	node2 := list2.head

	if node1 == nil && node2 == nil {
		return merged, nil
	}

	// для скорости формируем список 'руками', Add() дал бы O(N*N)
	var node *Node[T]
	addToMerged := func(val T) {
		newNode := &Node[T]{value: val, prev: node}

		if node == nil {
			merged.head = newNode
		} else {
			node.next = newNode
		}

		node = newNode
		merged.count++
	}

	for node1 != nil || node2 != nil {
		switch {
		case node1 == nil:
			addToMerged(node2.value)
			node2 = node2.next
		case node2 == nil:
			addToMerged(node1.value)
			node1 = node1.next
		case list1.Compare(node1.value, node2.value) < 0:
			addToMerged(node1.value)
			node1 = node1.next
		default:
			addToMerged(node2.value)
			node2 = node2.next
		}
	}
	merged.tail = node

	return merged, nil
}
