package main

import (
	"errors"
)

var ErrNegativeCacheSize = errors.New("cache size must be non-negative")

const indexStep = 1

const (
	stateEmpty slotState = iota
	stateOccupied
	stateDeleted
)

type slotState uint8

type NativeCache[T any] struct {
	hits       []int
	size       int
	slots      []string // keys
	slotsState []slotState
	values     []T
}

func NewNativeCache[T any](size int) (*NativeCache[T], error) {
	if size < 0 {
		return nil, ErrNegativeCacheSize
	}

	return &NativeCache[T]{
		hits:       make([]int, size),
		size:       size,
		slots:      make([]string, size),
		slotsState: make([]slotState, size),
		values:     make([]T, size),
	}, nil
}

func (nc *NativeCache[T]) hashFun(key string) int {
	var hash uint32

	for i := range len(key) {
		hash = hash*31 + uint32(key[i])
	}

	return int(hash&0x7FFFFFFF) % nc.size
}

func (nc *NativeCache[T]) incIdx(idx int) int {
	return (idx + indexStep) % nc.size
}

func (nc *NativeCache[T]) Get(key string) (T, bool) {
	if nc.size == 0 {
		var zero T

		return zero, false
	}

	idx := nc.hashFun(key)
	for range nc.size {
		if nc.slotsState[idx] == stateEmpty {
			var zero T

			return zero, false
		}

		if nc.slotsState[idx] == stateOccupied && nc.slots[idx] == key {
			nc.hits[idx]++

			return nc.values[idx], true
		}

		idx = nc.incIdx(idx)
	}

	var zero T
	return zero, false
}

func (nc *NativeCache[T]) Put(key string, value T) {
	if nc.size == 0 {
		return
	}

	idx := nc.hashFun(key)
	firstDeleted := -1

	for range nc.size {
		if nc.slotsState[idx] == stateDeleted && firstDeleted == -1 {
			firstDeleted = idx
		}

		if nc.slotsState[idx] == stateOccupied && nc.slots[idx] == key {
			nc.values[idx] = value

			return
		}

		if nc.slotsState[idx] == stateEmpty {
			break
		}

		idx = nc.incIdx(idx)
	}

	if firstDeleted != -1 {
		nc.insertSlot(firstDeleted, key, value)

		return
	}

	if nc.slotsState[idx] == stateEmpty {
		nc.insertSlot(idx, key, value)

		return
	}

	minIdx := nc.findMinIdx()
	nc.deleteSlot(minIdx)
	nc.insertSlot(minIdx, key, value)
}

func (nc *NativeCache[T]) insertSlot(idx int, key string, value T) {
	nc.slots[idx] = key
	nc.slotsState[idx] = stateOccupied
	nc.values[idx] = value
}

func (nc *NativeCache[T]) deleteSlot(idx int) {
	var zero T

	nc.slotsState[idx] = stateDeleted
	nc.slots[idx] = ""
	nc.values[idx] = zero
	nc.hits[idx] = 0
}

func (nc *NativeCache[T]) findMinIdx() int {
	minIdx := -1

	for i := range nc.size {
		if nc.slotsState[i] != stateOccupied {
			continue
		}

		if minIdx == -1 || nc.hits[i] < nc.hits[minIdx] {
			minIdx = i
		}
	}

	return minIdx
}
