package main

import (
	"constraints"
	"errors"
	"math/rand"
	"strings"
)

// Порядковый номер задания на курсе: 7
// Номер задачи из задания: 8
// Краткое название: метод удаления всех дубликатов из упорядоченного списка
// Time complexity: O(N)
// Memory complexity: O(1)
var ErrOrderingMismatch = errors.New("ordering must match for both lists")

func (l *OrderedList[T]) RemoveDuplicates() {
	if l.Count() < 2 {
		return
	}

	node := l.head
	next := node.next
	for next != nil {
		if l.Compare(next.value, node.value) == 0 {
			node.next = next.next

			if next.next != nil {
				next.next.prev = node
			} else {
				l.tail = node
			}

			l.count--

			next = node.next
		} else {
			node = node.next
			next = node.next
		}
	}
}

// Номер задачи из задания: 9
// Краткое название: алгоритм слияния двух упорядоченных списков в один, сохраняя порядок элементов
// Time complexity: O(N)
// Memory complexity: O(N)
func MergeOrderedLists[T constraints.Ordered](list1, list2 *OrderedList[T]) (*OrderedList[T], error) {
	if list1.ascending != list2.ascending {
		return &OrderedList[T]{}, ErrOrderingMismatch
	}

	merged := &OrderedList[T]{ascending: list1.ascending}

	node1 := list1.head
	node2 := list2.head

	if node1 == nil && node2 == nil {
		return merged, nil
	}

	// для скорости формируем список 'руками', Add() дал бы O(N*N)
	var node *Node[T]
	addToMerged := func(val T) {
		newNode := &Node[T]{value: val, prev: node}

		if node == nil {
			merged.head = newNode
		} else {
			node.next = newNode
		}

		node = newNode
		merged.count++
	}

	for node1 != nil || node2 != nil {
		switch {
		case node1 == nil:
			addToMerged(node2.value)
			node2 = node2.next
		case node2 == nil:
			addToMerged(node1.value)
			node1 = node1.next
		case list1.Compare(node1.value, node2.value) < 0:
			addToMerged(node1.value)
			node1 = node1.next
		default:
			addToMerged(node2.value)
			node2 = node2.next
		}
	}
	merged.tail = node

	return merged, nil
}

// Номер задачи из задания: 10
// Краткое название: метод проверки наличия заданного упорядоченного подсписка в текущем списке
// Time complexity: O(N)
// Memory complexity: O(1)
func (l *OrderedList[T]) ContainsSublist(sublist *OrderedList[T]) bool {
	if sublist.Count() == 0 {
		return true
	}
	if l.Count() < sublist.Count() {
		return false
	}
	if l.ascending != sublist.ascending {
		return false
	}

	node := l.head
	nodeCountLeft := l.Count()
	sub := sublist.head
	for node != nil {
		compare := l.Compare(node.value, sub.value)
		// текущий элемент l уже больше (или меньше) начала sublist
		if (l.ascending && compare > 0) || (!l.ascending && compare < 0) {
			break
		}
		if nodeCountLeft < sublist.Count() {
			return false
		}

		if compare == 0 {
			n := node
			s := sub
			match := true

			for s != nil {
				if n == nil || l.Compare(n.value, s.value) != 0 {
					match = false
					break
				}

				n = n.next
				s = s.next
			}

			if match {
				return true
			}
		}

		node = node.next
		nodeCountLeft--
	}

	return false
}

// Номер задачи из задания: 11
// Краткое название: метод, который находит наиболее часто встречающееся значение в списке
// Time complexity: O(N)
// Memory complexity: O(1)
func (l *OrderedList[T]) MostFrequentVal() (T, error) {
	if l.Count() == 0 {
		var zero T
		return zero, ErrEmptyList
	}

	node := l.head
	mostFreqVal := node.value
	maxCount := 1
	currentCount := 1
	currentVal := node.value

	for node := l.head.next; node != nil; node = node.next {
		if l.Compare(node.value, currentVal) == 0 {
			currentCount++
		} else {
			if currentCount > maxCount {
				maxCount = currentCount
				mostFreqVal = currentVal
			}

			currentVal = node.value
			currentCount = 1
		}
	}

	if currentCount > maxCount {
		mostFreqVal = currentVal
	}

	return mostFreqVal, nil
}

// Номер задачи из задания: 12
// Краткое название: возможность найти индекс элемента в списке, которая должна работать за O(log N)
var ErrInvalidIndex = errors.New("invalid index")

const maxLevel = 16 // log2(N) = 16, хранит до 65536 элементов

type SkipNode[T constraints.Ordered] struct {
	value T
	next  []*SkipNode[T] // next[i] - указатель на следущий узел на уровне i
	// span[i] - сколько элементов базового уровня мы перепрыгнем, переходя по указателю next[i].
	// Важно, что span[i] - это свойство указателя (ребра), а не узла.
	// Он описывает "расстояние" от текущего узла до next[i]
	span []int
}

type SkipList[T constraints.Ordered] struct {
	head      *SkipNode[T]
	level     int // текущая максимальная высота
	count     int
	ascending bool
	maxLevel  int
}

func NewSkipList[T constraints.Ordered](ascending bool) *SkipList[T] {
	return &SkipList[T]{
		head: &SkipNode[T]{
			next: make([]*SkipNode[T], maxLevel),
			span: make([]int, maxLevel),
		},
		ascending: ascending,
		maxLevel:  maxLevel,
	}
}

// Add добавляет элемент по значению
// Time complexity: O(log N)
// Memory complexity: O(N)
func (l *SkipList[T]) Add(value T) {
	update, rank := l.findUpdateArray(value)
	newLevel := randomLevel()

	// нужно добавить новый уровень
	if newLevel > l.level {
		levelDiff := newLevel - l.level

		for range levelDiff {
			update = append(update, l.head)
			rank = append(rank, 0)
		}

		l.level = newLevel
	}

	newNode := &SkipNode[T]{
		value: value,
		next:  make([]*SkipNode[T], newLevel),
		span:  make([]int, newLevel),
	}
	countBeforeAdd := l.count

	for i := range newLevel {
		newNode.next[i] = update[i].next[i]

		if update[i].next[i] == nil {
			// если за новым узлом на этом уровне ничего нет (хвост),
			// новый узел должен покрывать все элементы, оставшиеся до конца списка
			newNode.span[i] = (countBeforeAdd + 1) - rank[0]
		} else {
			// если за ним есть узел, забираем остаток span от предыдущего ребра
			newNode.span[i] = update[i].span[i] - (rank[0] - rank[i])
		}

		update[i].next[i] = newNode
		// новое ребро от update[i] до newNode покрывает расстояние между ними + 1 шаг на сам newNode
		update[i].span[i] = (rank[0] - rank[i]) + 1
	}

	// между update[i] и его next добавлен еще один узел -> увеличиваем span
	for i := newLevel; i < l.level; i++ {
		update[i].span[i]++
	}

	l.count++
}

// FindByIndex находит элемент по индексу
// Time complexity: O(log N)
// Memory complexity: O(1)
func (l *SkipList[T]) FindByIndex(index int) (T, error) {
	if index < 0 || index >= l.count {
		var zero T
		return zero, ErrInvalidIndex
	}

	node := l.head
	currIndex := 0

	for level := l.level - 1; level >= 0; level-- {
		// currIndex + node.span[level] - длина перехода от HEAD до node
		for node.next[level] != nil && currIndex+node.span[level] <= index+1 {
			currIndex += node.span[level]
			node = node.next[level]
		}
	}

	return node.value, nil
}

// Delete удаляет первое вхождение значения. Возвращает true, если элемент был найден и удалён.
// Time complexity: O(log N)
// Memory complexity: O(1)
func (l *SkipList[T]) Delete(value T) bool {
	update, _ := l.findUpdateArray(value)

	// update[0] - это узел, который находится строго перед тем местом, где должен находиться искомый value
	target := update[0].next[0]

	if target == nil || l.Compare(target.value, value) != 0 {
		return false
	}

	targetLevel := len(target.next)

	// удаление узла
	for i := range targetLevel {
		update[i].next[i] = target.next[i]
		// уменьшаем span ребра от предыдущего до следующего за удаляемым на 1 (минус удаляемый узел)
		update[i].span[i] += target.span[i] - 1
	}

	// у уровней, где удаляемого узла не было, но он находился в промежутке,
	// нужно уменьшить span на 1
	for i := targetLevel; i < l.level; i++ {
		update[i].span[i]--
	}

	l.count--

	// уменьшаем l.level, если верхние уровни опустели
	for l.level > 1 && l.head.next[l.level-1] == nil {
		l.level--
	}

	return true
}

func randomLevel() int {
	level := 1

	for level < maxLevel && rand.Float64() < 0.5 {
		level++
	}

	return level
}

// findUpdateArray возвращает массив узлов, после которых нужно вставить/удалить элемент
// на каждом уровне от 0 до l.level-1
func (l *SkipList[T]) findUpdateArray(value T) (update []*SkipNode[T], rank []int) {
	update = make([]*SkipNode[T], l.level) // l.level - кол-во активных уровней
	rank = make([]int, l.level)            // rank[i] хранит индекс узла update[i]

	node := l.head
	for lev := l.level - 1; lev >= 0; lev-- {
		// наследуем расстояние от HEAD до узла node
		if lev < l.level-1 { // спустились с lev+1
			rank[lev] = rank[lev+1]
		} else { // только вышли из HEAD (lev := l.level - 1)
			rank[lev] = 0
		}

		for node.next[lev] != nil {
			compare := l.Compare(node.next[lev].value, value)

			if !l.ascending {
				compare = -compare
			}

			if compare < 0 {
				rank[lev] += node.span[lev]
				node = node.next[lev]
			} else {
				break
			}
		}

		update[lev] = node
	}

	return update, rank
}

func (l *SkipList[T]) Compare(v1, v2 T) int {
	result := 0

	if s1, ok := any(v1).(string); ok {
		s2 := any(v2).(string)
		s1, s2 = strings.TrimSpace(s1), strings.TrimSpace(s2)

		if s1 < s2 {
			result = -1
		} else if s1 > s2 {
			result = +1
		}

		return result
	}

	if v1 < v2 {
		result = -1
	} else if v1 > v2 {
		result = +1
	}

	return result
}


