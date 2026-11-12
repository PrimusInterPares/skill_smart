package main

import (
	"constraints"
	"fmt"
	"os"
	"strconv"
	"hash/fnv"
)

const TableSize = 20000

type Element[T constraints.Ordered] struct {
	value    T
	occupied bool
}

type PowerSet[T constraints.Ordered] struct {
	elements []Element[T]
	count    int
}

func Init[T constraints.Ordered]() PowerSet[T] {
	return PowerSet[T]{
		elements: make([]Element[T], TableSize),
		count:    0,
	}
}

func (ps *PowerSet[T]) Size() int {
	return ps.count
}

func (ps *PowerSet[T]) isFull() bool {
	return ps.count == TableSize
}

func hashFun[T constraints.Ordered](value T, size int) int {
	str := fmt.Sprintf("%v", value)
	hash := fnv.New32a()
	hash.Write([]byte(str))

	return int(hash.Sum32() % uint32(size))
}

func (ps *PowerSet[T]) incIndex(index int) int {
	return (index + 1) % TableSize
}

func (ps *PowerSet[T]) Put(value T) {
	if ps.isFull() {
		return
	}

	index := hashFun(value, TableSize)

	for range TableSize {
		if ps.elements[index].occupied && ps.elements[index].value == value {
			return
		}

		if !ps.elements[index].occupied {
			ps.elements[index].value = value
			ps.elements[index].occupied = true
			ps.count++

			return
		}

		index = ps.incIndex(index)
	}
}

func (ps *PowerSet[T]) Get(value T) bool {
	return ps.findIndex(value) != -1
}

func (ps *PowerSet[T]) Remove(value T) bool {
	targetIndex := ps.findIndex(value)

	if targetIndex == -1 {
		return false
	}

	ps.elements[targetIndex].occupied = false
	ps.count--

	index := ps.incIndex(targetIndex)
	for ps.elements[index].occupied {
		v := ps.elements[index].value
		ps.elements[index].occupied = false
		ps.count--

		ps.Put(v)

		index = ps.incIndex(index)
	}

	return true
}

func (ps *PowerSet[T]) findIndex(value T) int {
	index := hashFun(value, TableSize)

	for range TableSize {
		if !ps.elements[index].occupied {
			return -1
		}

		if ps.elements[index].value == value {
			return index
		}

		index = ps.incIndex(index)
	}

	return -1
}

func (ps *PowerSet[T]) Intersection(set2 PowerSet[T]) PowerSet[T] {
	result := Init[T]()

	for _, el := range ps.elements {
		if el.occupied && set2.Get(el.value) {
			result.Put(el.value)
		}
	}

	return result
}

func (ps *PowerSet[T]) Union(set2 PowerSet[T]) PowerSet[T] {
	result := Init[T]()

	result.fillWith(ps)
	result.fillWith(&set2)

	return result
}

func (ps *PowerSet[T]) fillWith(set2 *PowerSet[T]) {
	for _, el := range set2.elements {
		if el.occupied {
			ps.Put(el.value)
		}
	}
}

func (ps *PowerSet[T]) Difference(set2 PowerSet[T]) PowerSet[T] {
	result := Init[T]()

	for _, el := range ps.elements {
		if el.occupied && !set2.Get(el.value) {
			result.Put(el.value)
		}
	}

	return result
}

func (ps *PowerSet[T]) IsSubset(set2 PowerSet[T]) bool {
	for _, el := range set2.elements {
		if el.occupied && !ps.Get(el.value) {
			return false
		}
	}

	return true
}

func (ps *PowerSet[T]) Equals(set2 PowerSet[T]) bool {
	return ps.Size() == set2.Size() && ps.IsSubset(set2)
}


