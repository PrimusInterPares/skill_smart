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
}

type Deque[T any] struct {
	front *Block[T]
	back  *Block[T]
	size  int
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
		d.front.addFront(itm)
		d.size++

		return
	}

	// front block is full (head == 0)
	b := newFrontBlock[T]()
	b.head--
	b.items[b.head] = itm

	b.next = d.front
	d.front.prev = b
	d.front = b
	d.size++
}

func newFrontBlock[T any]() *Block[T] {
	return newBlock[T](blockSize, blockSize)
}

func (b *Block[T]) addFront(itm T) {
	b.head--
	b.items[b.head] = itm
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
