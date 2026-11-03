package main

import (
	"errors"
)

var ErrNoLessFunc = errors.New("less function is not set, cannot calculate Min()")

func IsPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}

	deque := &Deque[rune]{}
	for _, r := range s {
		deque.AddFront(r)
	}

	for deque.Size() > 1 {
		left, _ := deque.RemoveFront()
		right, _ := deque.RemoveTail()

		if left != right {
			return false
		}
	}

	return true
}

func (d *Deque[T]) Min() (T, error) {
	var result T
	if d.size == 0 {
		return result, ErrEmptyDeque
	}
	if d.less == nil {
		return result, ErrNoLessFunc
	}

	var globalMin T
	found := false

	curr := d.front
	for curr != nil {
		for _, val := range curr.mins {
			if !found || d.less(val, globalMin) {
				globalMin = val
				found = true
			}
		}
		curr = curr.next
	}

	if !found {
		return result, ErrEmptyDeque
	}

	return globalMin, nil
}

func (d *Deque[T]) SetLess(less func(a, b T) bool) {
	d.less = less
}
