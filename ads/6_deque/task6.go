package main

import (
	"errors"
	// "os"
)

// Deque (size: 5)
//  │
//  ├── front ──┐
//  │           ▼
//  │    ┌───────────────────────┐       ┌───────────────────────┐
//  │    │ Block 1               │       │ Block 2               │
//  │    │ head: 1, tail: 4      │       │ head: 0, tail: 2      │
//  │    │ prev: nil             │ next  │ prev: &Block1         │
//  │    │                       ├──────►│ next: nil             │
//  │    │ data: [ _, A, B, C ]  │       │ data: [ D, E, _, _ ]  │
//  │    │          ▲        ▲   │       │        ▲        ▲     │
//  │    │        head      tail │       │      head      tail   │
//  │    └───────────────────────┘       └───────────────────────┘
//  │
//  └── back ─────────────────────────────────────────────────────┘

var ErrEmptyDeque = errors.New("Deque is empty")

const blockSize = 64

type Block[T any] struct {
	items [blockSize]T // slice of fixed size
	head  int          // first element index
	tail  int          // after last element index
	prev  *Block[T]
	next  *Block[T]
	mins  []T
}

type Deque[T any] struct {
	front *Block[T]
	back  *Block[T]
	size  int
	less  func(a, b T) bool
}

func (d *Deque[T]) Size() int {
	return d.size
}

func (d *Deque[T]) AddFront(itm T) {
	if d.front == nil {
		b := newFrontBlock[T]()
		d.front = b
		d.back = b
	} else if d.size == 0 {
		d.front.head = blockSize
		d.front.tail = blockSize
	}

	// front block has free space
	if d.front.head > 0 {
		d.front.addFront(itm, d.less)
		d.size++

		return
	}

	// front block is full (head == 0)
	b := newFrontBlock[T]()
	b.head--
	b.items[b.head] = itm

	if d.less != nil {
		b.mins = append(b.mins, itm)
	}

	b.next = d.front
	d.front.prev = b
	d.front = b
	d.size++
}

func newFrontBlock[T any]() *Block[T] {
	return newBlock[T](blockSize, blockSize)
}

func (b *Block[T]) addFront(itm T, less func(T, T) bool) {
	b.head--
	b.items[b.head] = itm

	if less != nil {
		for len(b.mins) > 0 && less(itm, b.mins[0]) {
			b.mins = b.mins[1:]
		}
		b.mins = append([]T{itm}, b.mins...)
	}
}

func (d *Deque[T]) AddTail(itm T) {
	if d.back == nil {
		b := newBackBlock[T]()
		d.front = b
		d.back = b
	} else if d.size == 0 {
		d.back.head = 0
		d.back.tail = 0
	}

	if d.less != nil {
		for len(d.back.mins) > 0 && d.less(itm, d.back.mins[len(d.back.mins)-1]) {
			d.back.mins = d.back.mins[:len(d.back.mins)-1]
		}
		d.back.mins = append(d.back.mins, itm)
	}

	// back block has free space
	if d.back.tail < blockSize {
		d.back.addTail(itm)
		d.size++

		return
	}

	// back block is full (tail == 64)
	b := newBackBlock[T]()
	b.items[b.head] = itm
	b.tail++

	b.prev = d.back
	d.back.next = b
	d.back = b
	d.size++
}

func newBackBlock[T any]() *Block[T] {
	return newBlock[T](0, 0)
}

func (b *Block[T]) addTail(itm T) {
	b.items[b.tail] = itm
	b.tail++
}

func (d *Deque[T]) RemoveFront() (T, error) {
	var result T
	if d.Size() == 0 {
		return result, ErrEmptyDeque
	}

	var zero T
	b := d.front
	result = b.items[b.head]
	b.items[b.head] = zero
	b.head++
	d.size--

	if d.less != nil && len(b.mins) > 0 && isEqual(result, b.mins[0], d.less) {
		b.mins = b.mins[1:]
	}

	// if block is empty now
	if b.tail == b.head {
		// and it is the last block
		if d.back == b {
			return result, nil
		}

		d.front.next.prev = nil
		d.front = d.front.next

		return result, nil
	}

	return result, nil
}

func (d *Deque[T]) RemoveTail() (T, error) {
	var result T
	if d.Size() == 0 {
		return result, ErrEmptyDeque
	}

	var zero T
	b := d.back
	b.tail--
	result = b.items[b.tail]
	b.items[b.tail] = zero
	d.size--

	if d.less != nil && len(b.mins) > 0 && isEqual(result, b.mins[len(b.mins)-1], d.less) {
		b.mins = b.mins[:len(b.mins)-1]
	}

	// if block is empty now
	if b.tail == b.head {
		// and it is the last block
		if d.front == b {
			return result, nil
		}

		d.back.prev.next = nil
		d.back = d.back.prev

		return result, nil
	}

	return result, nil
}

func newBlock[T any](head, tail int) *Block[T] {
	return &Block[T]{head: head, tail: tail}
}

func isEqual[T any](a, b T, less func(T, T) bool) bool {
	return !less(a, b) && !less(b, a)
}
