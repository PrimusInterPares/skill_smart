package main

import (
	"errors"
)

var ErrNoLessFunc = errors.New("less function is not set, cannot calculate Min()")

// Порядковый номер задания на курсе: 6
// Номер задачи из задания: 4
// Краткое название: функция, которая с помощью Deque проверяет, является ли некоторая строка палиндромом
func IsPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}

	deque := &Deque[rune]{}
	for _, r := range s {
		deque.AddFront(r)
	}

	for deque.Size() > 1 {
		left, _ := deque.RemoveFront()
		right, _ := deque.RemoveTail()

		if left != right {
			return false
		}
	}

	return true
}

// Номер задачи из задания: 5
// Краткое название: метод, который возвращает минимальный элемент деки за O(1)
// Time complexity: O(1)
// Memory complexity: O(1)
func (d *Deque[T]) Min() (T, error) {
	var result T

	if d.size == 0 {
		return result, ErrEmptyDeque
	}

	if d.less == nil {
		return result, ErrNoLessFunc
	}

	return d.minBuf[d.minHead], nil
}

func (d *Deque[T]) SetLess(less func(a, b T) bool) {
	d.less = less
}

// Номер задачи из задания: 6
// Краткое название: двусторонняя очередь с помощью динамического массива
const decreaseBufferRatio float64 = 1.5
const increaseBufferRatio int = 2
const minimalBufferSize int = 16

type DynArray[T any] struct {
	count    int
	capacity int
	array    []T
}

type DynArrayDeque[T any] struct {
	items *DynArray[T]
	head  int // index of the first element
	tail  int // index of the last element
}

func (d *DynArrayDeque[T]) Init() {
	d.items.count = 0
	d.resize(minimalBufferSize)
}

func (d *DynArrayDeque[T]) resize(sz int) {
	newArr := make([]T, sz)

	if d.items.count != 0 {
		idx := d.head
		for i := range d.items.count {
			newArr[(sz-d.items.count)/2+i] = d.items.array[(idx+i+d.items.capacity)%d.items.capacity]
		}
	}

	d.head = (sz - d.items.count) / 2
	d.tail = d.head + d.items.count - 1
	d.items.array = newArr
	d.items.capacity = sz
}

// Time complexity: O(1) until buffer is resized
// Memory complexity: O(1) until buffer is resized
func (d *DynArrayDeque[T]) AddFront(val T) {
	if d.items.count == d.items.capacity {
		d.increaseBuffer()
	}

	d.head = (d.head - 1 + d.items.capacity) % d.items.capacity
	d.items.array[d.head] = val
	d.items.count++
}

// Time complexity: O(1) until buffer is resized
// Memory complexity: O(1) until buffer is resized
func (d *DynArrayDeque[T]) RemoveHead() (T, error) {
	var result T

	if d.items.count == 0 {
		return result, ErrEmptyDeque
	}

	result = d.items.array[d.head]
	var zero T
	d.items.array[d.head] = zero // for GC
	d.head = (d.head + 1 + d.items.capacity) % d.items.capacity
	d.items.count--

	if d.items.count*2 < d.items.capacity {
		d.decreaseBuffer()
	}

	return result, nil
}

// Time complexity: O(1) until buffer is resized
// Memory complexity: O(1) until buffer is resized
func (d *DynArrayDeque[T]) AddTail(val T) {
	if d.items.count == d.items.capacity {
		d.increaseBuffer()
	}

	d.tail = (d.tail + 1) % d.items.capacity
	d.items.array[d.tail] = val
	d.items.count++
}

// Time complexity: O(1) until buffer is resized
// Memory complexity: O(1) until buffer is resized
func (d *DynArrayDeque[T]) RemoveTail() (T, error) {
	var result T

	if d.items.count == 0 {
		return result, ErrEmptyDeque
	}

	result = d.items.array[d.tail]
	var zero T
	d.items.array[d.tail] = zero // for GC
	d.tail = (d.tail - 1 + d.items.capacity) % d.items.capacity
	d.items.count--

	if d.items.count*2 < d.items.capacity {
		d.decreaseBuffer()
	}

	return result, nil
}

func (d *DynArrayDeque[T]) increaseBuffer() {
	newCap := d.items.capacity * increaseBufferRatio

	if newCap == 0 {
		newCap = minimalBufferSize
	}

	d.resize(newCap)
}

func (d *DynArrayDeque[T]) decreaseBuffer() {
	newCap := int(float64(d.items.capacity) / decreaseBufferRatio)

	d.resize(max(newCap, minimalBufferSize))
}

// Рефлексия
// Баланс открывающих и закрывающих скобок
// Я так и реализовал, и даже по началу сделал словрь, но потом удалил, т.к. он выглядел избыточно.
// Забавно, что эту задачу мне в своё время дали на собеседовании.
//
// Текущий минимальный элемент в стеке за O(1)
// Я использовал для хранения минимумов slice, хотя сейчас стек выглядит гораздо более адекватным выбором.
// Где у меня в мышлении косяк?...
//
// Среднее значение всех элементов в стеке за O(1)
// Да, ровно так я и реализовал, как Вы написали.
//
// Постфиксная запись выражения
// Я учитвал, что из стека мы забираем элементы в обратном порядке, поэтому делал b +/* a.
// А вот лямбды это классная идея, попробую-ка я с ними переписать...


