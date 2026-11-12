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
