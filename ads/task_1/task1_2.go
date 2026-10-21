package main

import(
	"errors"
)

// Порядковый номер задания на курсе: 1
// Номер задачи из задания: 8
// Краткое название: соответственно суммировать в новый список элементы данных списков, если их длины равны.
// Сложность по времени: O(N)
// Сложность по памяти: O(N)
func SumLists(first *LinkedList, second *LinkedList) (*LinkedList, error) {
	if first.Count() != second.Count() {
		return &LinkedList{}, errors.New("lists have different sizes")
	}

	returnList := &LinkedList{}
	firstNode := first.head
	secondNode := second.head

	for firstNode != nil {
		returnList.AddInTail(Node{value: firstNode.value + secondNode.value})

		firstNode = firstNode.next
		secondNode = secondNode.next
	}

	return returnList, nil
}


