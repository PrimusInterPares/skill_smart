package main

import (
	"errors"
)

// Порядковый номер задания на курсе: 5
// Номер задачи из задания: 3
// Краткое название: функция, 'вращающая' очередь по кругу на N элементов
// Time complexity: O(N)
// Memory complexity: O(N)
func (q *Queue[T]) Rotate(n int) {
	if q.Size() <= 1 || n == 0 {
		return
	}

	normalN := n % q.Size() // 10 / 3 = 1 реальный круг
	if normalN < 0 {
		normalN += q.Size()
	}

	q.items = append(q.items[normalN:], q.items[:normalN]...)
}

// Номер задачи из задания: 4
// Краткое название: очередь с помощью двух стеков
var ErrEmptyStack = errors.New("stack is empty")

type node[T any] struct {
	next  *node[T]
	value T
}

type Stack[T any] struct {
	head *node[T]
	size int
}

func (st *Stack[T]) Size() int {
	return st.size
}

func (st *Stack[T]) Peek() (T, error) {
	var result T
	if st.Size() == 0 {
		return result, ErrEmptyStack
	}

	result = st.head.value

	return result, nil
}

func (st *Stack[T]) Pop() (T, error) {
	var result T
	if st.Size() == 0 {
		return result, ErrEmptyStack
	}
	result = st.head.value
	st.head = st.head.next
	st.size--

	return result, nil
}

func (st *Stack[T]) Push(item T) {
	st.head = &node[T]{value: item, next: st.head}
	st.size++
}

type StackedQueue[T any] struct {
	stackIn  *Stack[T]
	stackOut *Stack[T]
}

func NewStackedQueue[T any]() *StackedQueue[T] {
	return &StackedQueue[T]{
		stackIn:  &Stack[T]{},
		stackOut: &Stack[T]{},
	}
}

func (q *StackedQueue[T]) Size() int {
	return q.stackIn.Size() + q.stackOut.Size()
}

// Time complexity: O(1)
// Memory complexity: O(1)
func (q *StackedQueue[T]) Push(itm T) {
	q.stackIn.Push(itm)
}

// Time complexity: O(1), если stackOut не пуст, иначе O(N)
// Memory complexity: O(1)
func (q *StackedQueue[T]) Peek() (T, error) {
	var result T

	if q.Size() == 0 {
		return result, ErrEmptyQueue
	}

	if q.stackOut.Size() == 0 {
		q.moveInToOut()
	}

	result, err := q.stackOut.Peek()
	if err != nil {
		return result, err
	}

	return result, nil
}

// Time complexity: O(1), если stackOut не пуст, иначе O(N)
// Memory complexity: O(1)
func (q *StackedQueue[T]) Pop() (T, error) {
	var result T

	if q.Size() == 0 {
		return result, ErrEmptyQueue
	}

	if q.stackOut.Size() == 0 {
		q.moveInToOut()
	}

	result, err := q.stackOut.Pop()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (q *StackedQueue[T]) moveInToOut() {
	for q.stackIn.Size() > 0 {
		inVal, _ := q.stackIn.Pop()
		q.stackOut.Push(inVal)
	}
}

// Номер задачи из задания: 5
// Краткое название: функция, которая обращает все элементы в очереди в обратном порядке
// Time complexity: O()
// Memory complexity: O()
func (q *Queue[T]) Reverse() {

}

// Номер задачи из задания: 6
// Краткое название: круговая (циклическая буферная) очередь со статическим массивом фиксированного размера
type RingQueue[T any] struct {
	items  []T
	head   int
	length int
}

func NewRingQueue[T any](capacity int) *RingQueue[T] {
	return &RingQueue[T]{items: make([]T, capacity)}
}

// Рефлексия
// Номер задачи из задания: 6
// Краткое название: Динамический массив на основе банковского метода.
// Я решил усреднить стоимость операции до 3, хотя по расчётам выходило, что амортизационная стоимость удаления - 2.
// Списывать я при уменьшении в реаллокации ничего не стал, но было оущение, что постоянно копящиеся 'деньги'
// это скорее баг, нежели фича.

// Номер задачи из задания: 7
// Краткое название: Многомерный динамический массив.
// Это снова так просто и так гениально! Когда читаешь, кажется, что на поверхности такое решение.
// Но сам я не додумался, к сожалению, мутил рекурсию и вложенные дин. массивы =(
