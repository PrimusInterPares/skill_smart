package main

import(
	"slices"
)

// Порядковый номер задания на курсе: 2
// Номер задачи из задания: 9
// Краткое название: метод, который "переворачивает" порядок элементов в связном списке, меняя его на противоположный.
// Сложность по времени: O(N)
// Сложность по памяти: O(1)
func ReverseList(l *LinkedList2) *LinkedList2 {
	// Empty list
	if l.head == nil {
		return l
	}

	// One-element list
	if l.Count() == 1 {
		return l
	}

	tail := l.tail
	l.tail = l.head
	node := l.head
	for node != nil {
		next := node.next

		node.prev, node.next = node.next, node.prev

		node = next
	}
	l.head = tail

	return l
}

// Номер задачи из задания: 10
// Краткое название: булев метод, который сообщает, имеются ли циклы (замкнутые на себя по кругу) внутри списка.
// Сложность по времени: O(N)
// Сложность по памяти: O(N)
func HasCycle(l *LinkedList2) bool {
	// Empty list
	if l.head == nil {
		return false
	}

	visited := make(map[*Node]struct{})
	node := l.head
	var prevNode *Node

	for node != nil {
		// forward
		if node.next != nil {
			if _, exists := visited[node.next]; exists {
				return true
			}
		}

		// backward
		if node.prev != nil && node.prev != prevNode {
			return true
		}

		visited[node] = struct{}{}
		prevNode = node
		node = node.next
	}

	return false
}

// Номер задачи из задания: 11
// Краткое название: метод, сортирующий список.
// Сложность по времени: O(N*logN)
// Сложность по памяти: O(N)
func SortList(l *LinkedList2) *LinkedList2 {
	// Empty or one-element list
	if l.head == nil || l.head.next == nil {
		return l
	}

	values := make([]int, 0, capacity(l))
	for node := l.head; node != nil; node = node.next {
		values = append(values, node.value)
	}
	slices.Sort(values)

	sorted := &LinkedList2{}
	for i, val := range values {
		newNode := Node{value: val}

		if i == 0 {
			sorted.InsertFirst(newNode)
		} else {
			sorted.Insert(sorted.tail, newNode)
		}
	}

	return sorted
}


func capacity(l *LinkedList2) int {
	count := 0
	for node := l.head; node != nil; node = node.next {
		count++
	}

	return count
}

// Номер задачи из задания: 12
// Краткое название: метод, сортирующий два списка и объединяющий их в отсортированный третий.
// Сложность по времени: O()
// Сложность по памяти: O()
func MergeLists(first *LinkedList2, second *LinkedList2) *LinkedList2 {
	return &LinkedList2{}
}


