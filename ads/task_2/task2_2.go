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
func MergeLists(first, second *LinkedList2) *LinkedList2 {
	if first.head == nil {
		return SortList(second)
	}

	if second.head == nil {
		return SortList(first)
	}
	
	sortedFirst := SortList(first)
	sortedSecond := SortList(second)
	merged := &LinkedList2{}

	var mergeFirst *LinkedList2
	var mergeSecond *LinkedList2
	// First element of merged will be <= any element from any sorted list merged in second
	if sortedFirst.head.value < sortedSecond.head.value {
		mergeFirst = sortedFirst
		mergeSecond = sortedSecond
	} else {
		mergeFirst = sortedSecond
		mergeSecond = sortedFirst
	}

	i := 0
	for node := mergeFirst.head; node != nil; node = node.next {
		newNode := Node{value: node.value}

		if i == 0 {
			merged.InsertFirst(newNode)
			i++
		} else {
			merged.Insert(merged.tail, newNode)
		}
	}

	for node := mergeSecond.head; node != nil; node = node.next {
		func() {
			for mNode := merged.head; mNode != nil; mNode = mNode.next {
				if node.value >= mNode.value {
					if mNode == merged.tail {
						merged.Insert(mNode, Node{value: node.value})
						return
					}
					
					mNodeNext := mNode.next
					if mNodeNext.value > node.value {
						merged.Insert(mNode, Node{value: node.value})
						return
					}
				}
			}
		}()
	}
	
	return merged
}


