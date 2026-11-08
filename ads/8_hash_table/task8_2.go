package main

import (
	"math"
)

// Порядковый номер задания на курсе: 8
// Номер задачи из задания: 3
// Краткое название: динамическая хэш-таблица, которая автоматически увеличивает свой размер при необходимости
const resizeLimit = 0.7

type DynHashTable struct {
	*HashTable
	count int // кол-во уникальных элементов
}

func InitDynHashTable(sz, stp int) *DynHashTable {
	baseTable := Init(sz, stp)

	return &DynHashTable{
		HashTable: &baseTable,
		count:     0,
	}
}

func (dht *DynHashTable) Put(value string) int {
	if dht.needsResize() {
		dht.resize()
	}

	index := dht.SeekSlot(value)

	if index == -1 {
		return -1
	}

	if dht.statuses[index] != statusOccupied {
		dht.count++
	}
	dht.slots[index] = value
	dht.statuses[index] = statusOccupied

	return index
}

func (dht *DynHashTable) resize() {
	oldSlots := dht.slots
	oldStatuses := dht.statuses

	newHt := Init(dht.newSize(), StepLength)
	dht.HashTable = &newHt

	for i := range len(oldSlots) {
		if oldStatuses[i] != statusOccupied {
			continue
		}

		value := oldSlots[i]

		// считаем хэш, избегая оверхед dht.Put
		index := dht.SeekSlot(value)
		dht.slots[index] = value
		dht.statuses[index] = statusOccupied
	}
}

// loadFactor возвращает значение от 0.0 до 1.0
func (dht *DynHashTable) loadFactor() float64 {
	return float64(dht.count) / float64(dht.size)
}

func (dht *DynHashTable) needsResize() bool {
	return dht.loadFactor() > resizeLimit
}

func (dht *DynHashTable) newSize() int {
	approxSize := dht.size * 2

	// берем ПЕРВОЕ простое число
	for i := approxSize + 1; ; i++ {
		if isPrime(i) {
			return i
		}
	}
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	limit := int(math.Sqrt(float64(n)))
	for i := 3; i <= limit; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// Номер задачи из задания: 4
// Краткое название: хэш-таблица, которая использует несколько хэш-функций для каждой операции вставки
type DoubleHashTable struct {
	*HashTable
}

func (dht *DoubleHashTable) SeekSlot(value string) int {
	index := dht.h1(value)
	step := dht.h2(value)
	firstTombstone := -1

	for range dht.size {
		status := dht.statuses[index]

		if status == statusEmpty && firstTombstone != -1 {
			return firstTombstone
		}

		if status == statusEmpty {
			return index
		}

		if status == statusOccupied && dht.slots[index] == value {
			return index
		}

		if status == statusDeleted && firstTombstone == -1 {
			firstTombstone = index
		}

		// на первой итерации цикла пробуем найти место для индекса из h1,
		// иначе подмешиваем h2
		index = (index + step) % dht.size
	}

	return firstTombstone
}

// полиномиальное хэширование: каждый символ строки в unsigned int прибавляется к хэшу, умноженному на 31 (как в Java)
func (dht *DoubleHashTable) h1(value string) int {
	var hash uint32

	for i := range len(value) {
		hash = hash*31 + uint32(value[i])
	}

	// маска 0x7FFFFFFF помоагет безопасно привести к int
	return int(hash&0x7FFFFFFF) % dht.size
}

// алгоритм FNV-1a
func (dht *DoubleHashTable) h2(value string) int {
	var hash uint32 = 2166136261 // offsetBasis
	var primeMultiplier uint32 = 16777619

	for i := range len(value) {
		hash = (hash ^ uint32(value[i])) * primeMultiplier
	}

	return int(hash)%(dht.size-1) + 1
}

// Рефлексия
// Для минимума в деке я добавил кольцевой стек. Сильно это хуже второй деки? Минус - это разве что фиксированный размер. Надо подумать...
// Для двусторонней очереди использовал динамический массив, но в методах деки я взаимодействовал с динамическим массивом напрямую, не инкапсулируя его логику (ведь так быстрее). Учту, исправлюсь.
// Коммент по стилю кодирования принят. СИ прошерстил, написал линтер для LazyVim.


