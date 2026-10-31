package main

import (
	"errors"
	"os"
)

var ErrEmptyQueue = errors.New("queue is empty")

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Time complexity: O(1)
// Memory complexity: O(1) in Go, as it just reestablish the beginning of the slice.
// But if there was a copying on an array, O(N) is the complexity.
func (q *Queue[T]) Dequeue() (T, error) {
	var result T

	if q.Size() == 0 {
		return result, ErrEmptyQueue
	}

	result = q.items[0]
	var zero T
	q.items[0] = zero
	q.items = q.items[1:]

	return result, nil
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (q *Queue[T]) Enqueue(itm T) {
	q.items = append(q.items, itm)
}


