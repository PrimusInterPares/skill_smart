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

	if lld.head.key == key {
		lld.head.value = value

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

// Номер задачи из задания: 6
// Краткое название: словарь, в котором ключи представлены битовыми строками фиксированной длины
type BinaryNode[T any] struct {
	next     [2]*BinaryNode[T] // индекс 0 заменяет left, индекс 1 заменяет right
	hasValue bool
	value    T
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (bn *BinaryNode[T]) SetValue(value T) {
	bn.value = value
	bn.hasValue = true
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (bn *BinaryNode[T]) RemoveValue() {
	var zero T

	bn.value = zero
	bn.hasValue = false
}

type BinaryTrieDictionary[T any] struct {
	root   *BinaryNode[T]
	bitLen int
	count  int
}

func NewBinaryTrieDictionary[T any](bitLen int) *BinaryTrieDictionary[T] {
	return &BinaryTrieDictionary[T]{
		root:   &BinaryNode[T]{},
		bitLen: bitLen,
		count:  0,
	}
}

// Time complexity: O(bitLen)
// Memory complexity: O(1)
func (btd *BinaryTrieDictionary[T]) IsKey(key uint32) bool {
	node := btd.root

	for i := btd.bitLen - 1; i >= 0; i-- {
		bit := (key >> i) & 1

		if node.next[bit] == nil {
			return false
		}

		node = node.next[bit]
	}

	return node.hasValue
}

// Time complexity: O(bitLen)
// Memory complexity: O(1)
func (btd *BinaryTrieDictionary[T]) Get(key uint32) (T, error) {
	var zero T
	node := btd.root

	for i := btd.bitLen - 1; i >= 0; i-- {
		bit := (key >> i) & 1

		if node.next[bit] == nil {
			return zero, ErrKeyNotFound
		}

		node = node.next[bit]
	}

	if !node.hasValue {
		return zero, ErrKeyNotFound
	}

	return node.value, nil
}

// Time complexity: O(bitLen)
// Memory complexity: O(bitLen) в худшем случае при аллокации новой ветки
func (btd *BinaryTrieDictionary[T]) Put(key uint32, value T) {
	node := btd.root

	for i := btd.bitLen - 1; i >= 0; i-- {
		bit := (key >> i) & 1

		if node.next[bit] == nil {
			node.next[bit] = &BinaryNode[T]{}
		}

		node = node.next[bit]
	}

	// новый узел
	if !node.hasValue {
		btd.count++
	}

	node.SetValue(value)
}

// Time complexity: O(bitLen)
// Memory complexity: O(1)
func (btd *BinaryTrieDictionary[T]) Delete(key uint32) error {
	node := btd.root

	for i := btd.bitLen - 1; i >= 0; i-- {
		bit := (key >> i) & 1

		if node.next[bit] == nil {
			return ErrKeyNotFound
		}

		node = node.next[bit]
	}

	if !node.hasValue {
		return ErrKeyNotFound
	}

	node.RemoveValue()
	btd.count--

	return nil
}

// Рефлексия
// Я очень рад, что логику всю сделал правильно. Но с точки зрения стиля - грязно, сам прям смотрю на код и ругаюсь.
// Для поиска индекса за O(log N) я реализовал Skip List, который упоминался в теор. части. Судя по всему, чуть overengineered, но я очень благодарен этому опыту.


