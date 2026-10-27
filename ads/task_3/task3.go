package main

import (
	"os"
	"fmt"
)

const decreaseBufferRatio float64 = 1.5
const increaseBufferRatio int = 2
const minimalBufferSize int = 16

type DynArray[T any] struct {
	count int
	capacity int
	array []T
}

func (da *DynArray[T]) Init() {
  da.count = 0
  da.MakeArray(minimalBufferSize)
}

func (da *DynArray[T]) MakeArray(sz int) {
	newArr := make([]T, sz)
	copy(newArr, da.array)
  da.array = newArr
  da.capacity = sz
}

// Сложность по времени: O(1) при вставке в конец (O(N) при увеличении буфера), O(N) при вставке в i-ю позицию
// Сложность по памяти: сама вставка - O(1), но при изменении буфера и копировании массива это O(N)
func (da *DynArray[T]) Insert(itm T, index int) (error) {
  if index < 0 || index > da.count {
		return fmt.Errorf("bad index '%d'", index)
	}

	if da.count == da.capacity {
		da.increaseBuffer()
	}

	for i:= da.count; i >= index; i-- {
		da.array[i + 1] = da.array[i]
	}
	da.array[index] = itm
	da.count++

	return nil
}

// Сложность по времени: O(N), т.к. сдвигаем все элементы
// Сложность по памяти: само удаление - O(1), но при изменении буфера и копировании массива это O(N)
func (da *DynArray[T]) Remove(index int) (error) {
  if index < 0 || index >= da.count {
		return fmt.Errorf("bad index '%d'", index)
	}

	for i := index; i < da.count - 1; i++ {
		da.array[i] = da.array[i + 1]
	}

	var dummy T
	da.array[da.count - 1] = dummy

	da.count--

	if da.count * 2 < da.capacity {
		da.decreaseBuffer()
	}

	return nil
}

func (da *DynArray[T]) Append(itm T) {
	if da.count == da.capacity {
		da.increaseBuffer()
	}
	
	da.array[da.count] = itm
	da.count++
}

func (da *DynArray[T]) GetItem(index int) (T, error) {
	if index < 0 || index >= da.count {
		var dummy T
		return dummy, fmt.Errorf("bad index '%d'", index)
	}

	return da.array[index], nil
}

func (da *DynArray[T]) increaseBuffer() {
	newCap := da.capacity * increaseBufferRatio
	
	if newCap == 0 {
		newCap = minimalBufferSize
	}

	da.MakeArray(newCap)
}

func (da *DynArray[T]) decreaseBuffer() {
	newCap := int(float64(da.capacity) / decreaseBufferRatio)

	da.MakeArray(max(newCap, minimalBufferSize))
}


