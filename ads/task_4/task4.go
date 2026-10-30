package main

import (
	"errors"
	"os"
)

var _ = os.Args
var ErrEmptyStack = errors.New("stack is empty")

type node[T any] struct {
	next  *node[T]
	value T
}

type Stack[T any] struct {
	head *node[T]
	size int
	less func(a, b T) bool
	mins []T
	sum  float64
}

func (st *Stack[T]) Size() int {
	return st.size
}

func (st *Stack[T]) Peek() (T, error) {
	var result T
	if st.Size() == 0 {
		return result, ErrEmptyStack
	}

	result = st.head.value

	return result, nil
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (st *Stack[T]) Pop() (T, error) {
	var result T
	if st.Size() == 0 {
		return result, ErrEmptyStack
	}
	result = st.head.value
	st.head = st.head.next
	st.size--

	// For Min()
	if len(st.mins) > 0 {
		st.mins = st.mins[:len(st.mins)-1]
	}
	// For Avg()
	st.sum -= toFloat64(result)
	return result, nil
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (st *Stack[T]) Push(item T) {
	st.head = &node[T]{value: item, next: st.head}
	st.size++

	// For Min()
	if st.less != nil {
		lastIndex := len(st.mins) - 1
		if len(st.mins) == 0 || st.less(item, st.mins[lastIndex]) {
			st.mins = append(st.mins, item)
		} else {
			st.mins = append(st.mins, st.mins[lastIndex])
		}
	}
	// For Avg()
	st.sum += toFloat64(item)
}

//--------------------------------------------------------------------------------------

type StackSlice[T any] struct {
	items []T
}

func (st *StackSlice[T]) Size() int {
	return len(st.items)
}

func (st *StackSlice[T]) Peek() (T, error) {
	var result T
	if len(st.items) == 0 {
		return result, ErrEmptyStack
	}

	result = st.items[st.Size()-1]

	return result, nil
}

// Time complexity: O(1) - average case, O(N) - if reallocation happens
// Memory complexity: O(1) - average case, O(N) - if reallocation happens
func (st *StackSlice[T]) Pop() (T, error) {
	var result T
	if len(st.items) == 0 {
		return result, ErrEmptyStack
	}

	topIndex := st.Size() - 1
	result = st.items[topIndex]

	var zero T
	st.items[topIndex] = zero // allow GC to collect it
	st.items = st.items[:topIndex]

	// Backing array must be resized 'manually'.
	// This code is needed if:
	// 1. there are harsh RAM limits,
	// 2. work with the stack will have peaks (stack raises to millions of elements and drops to almost empty).
	//
	// 1/4 is choosen to avoid `thrashing`: with 1/2 close Pop and Push will trigger expensive array copy too frequently.
	// 1024 is choosen because hanging RAM is cheaper for little stacks than CPU.
	stackCap := cap(st.items)
	if len(st.items) < stackCap/4 && stackCap > 1024 {
		st.items = append(make([]T, 0, stackCap/2), st.items...)
	}

	return result, nil
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (st *StackSlice[T]) Push(item T) {
	st.items = append(st.items, item)
}
