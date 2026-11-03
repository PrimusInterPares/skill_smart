package main

import (
	"errors"
	"os"
)

var ErrEmptyDeque = errors.New("Deque is empty")

const blockSize = 64

type Block[T any] struct {
	items [blockSize]T
	head  int
	tail  int
	prev  *Block[T]
	next  *Block[T]
}

type Deque[T any] struct {
	front *Block[T]
	back  *Block[T]
	size  int

	less    func(a, b T) bool
	minBuf  []T // ring stack
	minHead int
	minTail int
}

func (d *Deque[T]) Size() int {
	return d.size
}

func (d *Deque[T]) adjustMinCapacity() {
	if d.minBuf == nil {
		d.minBuf = make([]T, 128)
		d.minHead = 64
		d.minTail = 64

		return
	}

	if d.minHead == 0 || d.minTail == len(d.minBuf) {
		newCap := len(d.minBuf) * 2
		newBuf := make([]T, newCap)
		newHead := newCap / 4 // center
		n := copy(newBuf[newHead:], d.minBuf[d.minHead:d.minTail])

		d.minBuf = newBuf
		d.minHead = newHead
		d.minTail = newHead + n
	}
}

func (d *Deque[T]) pushMinFront(itm T) {
	if d.less == nil {
		return
	}

	d.adjustMinCapacity()

	for d.minHead < d.minTail && d.less(itm, d.minBuf[d.minHead]) {
		d.minHead++
	}
	d.minHead--
	d.minBuf[d.minHead] = itm
}

func (d *Deque[T]) pushMinTail(itm T) {
	if d.less == nil {
		return
	}

	d.adjustMinCapacity()

	for d.minHead < d.minTail && d.less(itm, d.minBuf[d.minTail-1]) {
		d.minTail--
	}
	d.minBuf[d.minTail] = itm
	d.minTail++
}

// Time complexity: O(1) amortized
// Memory complexity: O(1)
func (d *Deque[T]) AddFront(itm T) {
	if d.front == nil {
		d.front = &Block[T]{head: blockSize, tail: blockSize}
		d.back = d.front
	} else if d.size == 0 {
		d.front.head = blockSize
		d.front.tail = blockSize
	}

	if d.front.head > 0 {
		d.front.head--
		d.front.items[d.front.head] = itm
		d.size++

		d.pushMinFront(itm)

		return
	}

	b := &Block[T]{head: blockSize - 1, tail: blockSize}
	b.items[b.head] = itm
	b.next = d.front
	d.front.prev = b
	d.front = b
	d.size++

	d.pushMinFront(itm)
}

// Time complexity: O(1) amortized
// Memory complexity: O(1)
func (d *Deque[T]) AddTail(itm T) {
	if d.back == nil {
		d.back = &Block[T]{head: 0, tail: 0}
		d.front = d.back
	} else if d.size == 0 {
		d.back.head = 0
		d.back.tail = 0
	}

	if d.back.tail < blockSize {
		d.back.items[d.back.tail] = itm
		d.back.tail++
		d.size++

		d.pushMinTail(itm)

		return
	}

	b := &Block[T]{head: 0, tail: 1}
	b.items[0] = itm
	b.prev = d.back
	d.back.next = b
	d.back = b
	d.size++

	d.pushMinTail(itm)
}

// Time complexity: O(1) amortized
// Memory complexity: O(1)
func (d *Deque[T]) RemoveFront() (T, error) {
	var result T

	if d.size == 0 {
		return result, ErrEmptyDeque
	}

	var zero T
	b := d.front
	result = b.items[b.head]
	b.items[b.head] = zero
	b.head++
	d.size--

	if d.less != nil && d.minHead < d.minTail && isEqual(result, d.minBuf[d.minHead], d.less) {
		d.minHead++
	}

	if b.tail == b.head {
		if d.back == b {
			return result, nil // reuse last block
		}

		d.front = d.front.next
		d.front.prev = nil
	}

	return result, nil
}

// Time complexity: O(1) amortized
// Memory complexity: O(1)
func (d *Deque[T]) RemoveTail() (T, error) {
	var result T

	if d.size == 0 {
		return result, ErrEmptyDeque
	}

	var zero T
	b := d.back
	b.tail--
	result = b.items[b.tail]
	b.items[b.tail] = zero
	d.size--

	if d.less != nil && d.minHead < d.minTail && isEqual(result, d.minBuf[d.minTail-1], d.less) {
		d.minTail--
	}

	if b.tail == b.head {
		if d.front == b {
			return result, nil // reuse last block
		}

		d.back = d.back.prev
		d.back.next = nil
	}

	return result, nil
}

func isEqual[T any](a, b T, less func(T, T) bool) bool {
	return !less(a, b) && !less(b, a)
}


