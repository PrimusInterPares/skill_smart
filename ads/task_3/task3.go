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

// Time complexity: O(1) when iserted to tail (O(N) if bufer is increased), O(N) while iserting in i position
// Memory complexity: insert itself - O(1), but if buffer is resized array copy takes O(N)
func (da *DynArray[T]) Insert(item T, index int) (error) {
  if index < 0 || index > da.count {
		return fmt.Errorf("bad index '%d'", index)
	}

	if da.count == da.capacity {
		da.increaseBuffer()
	}

	for i:= da.count; i >= index; i-- {
		da.array[i + 1] = da.array[i]
	}
	da.array[index] = item
	da.count++

	return nil
}

// Time complexity: O(N) as we move i elements
// Memory complexity: remove itself - O(1), but if buffer is resized array copy takes O(N)
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

func (da *DynArray[T]) Append(item T) {
	if da.count == da.capacity {
		da.increaseBuffer()
	}
	
	da.array[da.count] = item
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


