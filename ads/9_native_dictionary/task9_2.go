package main

import (
	"constraints"
	"errors"
)

// Порядковый номер задания на курсе: 9
// Номер задачи из задания: 5
// Краткое название: словарь с использованием упорядоченного списка по ключу
var ErrEmptyDictionary = errors.New("dictionary is empty")

type Node[K constraints.Ordered, T any] struct {
	key   K
	value T
	next  *Node[K, T]
}

func (n *Node[K, T]) InsertAfter(key K, value T) *Node[K, T] {
	newNode := &Node[K, T]{
		key:   key,
		value: value,
		next:  n.next,
	}

	n.next = newNode

	return newNode
}

func (n *Node[K, T]) RemoveNext() *Node[K, T] {
	if n.next == nil {
		return nil
	}

	removed := n.next
	n.next = n.next.next

	return removed
}

type LinkedListDictionary[K constraints.Ordered, T any] struct {
	head  *Node[K, T]
	count int
}

func NewLinkedListDictionary[K constraints.Ordered, T any]() *LinkedListDictionary[K, T] {
	return &LinkedListDictionary[K, T]{
		head:  nil,
		count: 0,
	}
}

// Time complexity: O(N), хотя в среднем - быстрее
// Memory complexity: O(1)
func (lld *LinkedListDictionary[K, T]) IsKey(key K) bool {
	if lld.count == 0 {
		return false
	}

	if lld.head.key > key {
		return false
	}

	node := lld.head
	for node != nil {
		if node.key == key {
			return true
		}

		if node.key > key {
			return false
		}

		node = node.next
	}

	return false
}

// Time complexity: O(N), хотя в среднем - быстрее
// Memory complexity: O(1)
func (lld *LinkedListDictionary[K, T]) Get(key K) (T, error) {
	var zero T

	if lld.count == 0 {
		return zero, ErrEmptyDictionary
	}

	node := lld.head
	for node != nil {
		if node.key == key {
			return node.value, nil
		}

		if node.key > key {
			return zero, ErrKeyNotFound
		}

		node = node.next
	}

	return zero, ErrKeyNotFound
}

// Time complexity: O(N)
// Memory complexity: O(1)
func (lld *LinkedListDictionary[K, T]) Put(key K, value T) {
	if lld.count == 0 {
		lld.head = &Node[K, T]{
			key:   key,
			value: value,
			next:  nil,
		}

		lld.count++

		return
	}

	if lld.head.key > key {
		newNode := &Node[K, T]{
			key:   key,
			value: value,
			next:  lld.head,
		}

		lld.head = newNode
		lld.count++

		return
	}

	prev := lld.head
	node := prev.next
	for node != nil {
		if node.key == key {
			node.value = value

			return
		}

		if node.key > key {
			break
		}

		prev = node
		node = node.next
	}

	if lld.head.key == key {
		lld.head.value = value

		return
	}

	prev.InsertAfter(key, value)
	lld.count++
}

// Time complexity: O(N)
// Memory complexity: O(1)
func (lld *LinkedListDictionary[K, T]) Delete(key K) error {
	if lld.count == 0 {
		return ErrKeyNotFound
	}

	if lld.head.key == key {
		lld.head = lld.head.next
		lld.count--

		return nil
	}

	prev := lld.head
	node := prev.next
	for node != nil {
		if node.key == key {
			prev.RemoveNext()
			lld.count--

			return nil
		}

		if node.key > key {
			return ErrKeyNotFound
		}

		prev = node
		node = node.next
	}

	return ErrKeyNotFound
}
