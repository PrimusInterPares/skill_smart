package main

import (
	"constraints"
	"fmt"
)

// Порядковый номер задания на курсе: 10
// Номер задачи из задания: 4
// Краткое название: метод, реализующий декартово произведение множеств. Возвращает множество из пар соответствий.
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

// Номер задачи из задания: 6
// Краткое название: мульти-множество (Bag), в котором каждый элемент может присутствовать несколько раз
