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
