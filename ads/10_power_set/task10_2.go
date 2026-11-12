package main

import (
	"constraints"
	"fmt"
)

// Порядковый номер задания на курсе: 10
// Номер задачи из задания: 4
// Краткое название: метод, реализующий декартово произведение множеств. Возвращает множество из пар соответствий.
// Time complexity: O(N * M)
// Memory complexity: O(1)
func (ps *PowerSet[T]) CartesianProduct(set2 PowerSet[T]) PowerSet[string] {
	result := Init[string]()

	for _, el := range ps.elements {
		if el.occupied {
			addPairsFor(&result, el.value, &set2)
		}
	}

	return result
}

func addPairsFor[T constraints.Ordered](result *PowerSet[string], first T, set2 *PowerSet[T]) {
	for _, el := range set2.elements {
		if el.occupied {
			result.Put(fmt.Sprintf("(%v, %v)", first, el.value))
		}
	}
}

// Номер задачи из задания: 5
// Краткое название: функцию, которая находит пересечение любых трёх и более множеств (принимает количество множеств >= 3 в качестве списка).
// Time complexity: O(K * Mmin) или O(K * C) для фиксированного размера массива
// Memory complexity: O(K) или O(1)
func MultiIntersection[T constraints.Ordered](sets []PowerSet[T]) PowerSet[T] {
	if len(sets) == 0 {
		return Init[T]()
	}

	if len(sets) == 1 {
		return sets[0].Union(Init[T]())
	}

	minIndex := findMinIndex(sets)
	result := sets[minIndex].Union(Init[T]())
	for i, set := range sets {
		if i == minIndex {
			continue
		}

		result = result.Intersection(set)
	}

	return result
}

func findMinIndex[T constraints.Ordered](sets []PowerSet[T]) int {
	minSize := sets[0].Size()
	minIndex := 0

	for i := 1; i < len(sets); i++ {
		if sets[i].Size() < minSize {
			minSize = sets[i].Size()
			minIndex = i
		}
	}

	return minIndex
}

// Номер задачи из задания: 6
// Краткое название: мульти-множество (Bag), в котором каждый элемент может присутствовать несколько раз
type BagElement[T constraints.Ordered] struct {
	value T
	count int
}

func (be *BagElement[T]) isOccupied() bool {
	return be.count > 0
}

type Bag[T constraints.Ordered] struct {
	elements    []BagElement[T]
	uniqueCount int
}

func InitBag[T constraints.Ordered](size int) Bag[T] {
	if size <= 0 {
		size = TableSize
	}

	return Bag[T]{
		elements:    make([]BagElement[T], size),
		uniqueCount: 0,
	}
}

func (b *Bag[T]) Cap() int {
	return cap(b.elements)
}

func (b *Bag[T]) Size() int {
	return b.uniqueCount
}

func (b *Bag[T]) incIndex(index int) int {
	return (index + 1) % b.Cap()
}

func (b *Bag[T]) findIndex(value T) int {
	index := hashFun(value, b.Cap())

	for range b.Cap() {
		if !b.elements[index].isOccupied() {
			return -1
		}

		if b.elements[index].value == value {
			return index
		}

		index = b.incIndex(index)
	}

	return -1
}

func (b *Bag[T]) putWithCount(value T, count int) {
	index := hashFun(value, b.Cap())

	for range b.Cap() {
		if b.elements[index].isOccupied() && b.elements[index].value == value {
			b.elements[index].count += count

			return
		}

		if !b.elements[index].isOccupied() {
			b.elements[index].value = value
			b.elements[index].count += count
			b.uniqueCount++

			return
		}

		index = b.incIndex(index)
	}
}

func (b *Bag[T]) Put(value T) {
	b.putWithCount(value, 1)
}

func (b *Bag[T]) Remove(value T) {
	targetIndex := b.findIndex(value)

	if targetIndex == -1 {
		return
	}

	b.elements[targetIndex].count--
	if b.elements[targetIndex].isOccupied() {
		return
	}

	b.uniqueCount--
	index := b.incIndex(targetIndex)
	for b.elements[index].isOccupied() {
		v := b.elements[index].value
		c := b.elements[index].count
		b.elements[index].count = 0
		b.uniqueCount--

		b.putWithCount(v, c)

		index = b.incIndex(index)
	}
}

// Time complexity: O(K)
// Memory complexity O(K)
func (b *Bag[T]) GetElements() []BagElement[T] {
	result := make([]BagElement[T], 0, b.Size())

	for _, el := range b.elements {
		if el.isOccupied() {
			result = append(result, BagElement[T]{
				value: el.value,
				count: el.count,
			})
		}
	}

	return result
}

func findCountInSlice[T constraints.Ordered](elements []BagElement[T], value T) int {
	for _, el := range elements {
		if el.value == value {
			return el.count
		}
	}

	return 0
}

// Рефлексия
//
