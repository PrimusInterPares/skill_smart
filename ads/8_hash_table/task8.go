package main

import (
// "os"
// "strconv"
)

const (
	HashSize            = 19
	statusEmpty    byte = 0 // Slot is free
	statusOccupied byte = 1 // Slot is occupied
	statusDeleted  byte = 2 // Slot is deleted (tombstone)
	StepLength          = 3
)

type HashTable struct {
	size     int
	step     int
	slots    []string
	statuses []byte
}

func Init(sz, stp int) HashTable {
	ht := HashTable{
		size:     sz,
		step:     stp,
		slots:    make([]string, sz),
		statuses: make([]byte, sz),
	}

	return ht
}

func (ht *HashTable) HashFun(value string) int {
	var hash uint32 = 5381

	for i := range len(value) {
		// hash * 33 + c
		// (hash << 5) -> hash * 2^5 -> hash * 32
		// (hash * 32) + hash -> 33 * hash
		hash = ((hash << 5) + hash) + uint32(value[i])
	}

	// limit hash with HashTable.size
	// 0x7FFFFFFF mask guarantees the value will fit positive int
	return int(hash&0x7FFFFFFF) % ht.size
}

func (ht *HashTable) SeekSlot(value string) int {
	index := ht.HashFun(value)
	firstTombstone := -1

	for range ht.size {
		status := ht.statuses[index]

		if status == statusEmpty {
			if firstTombstone != -1 {
				index = firstTombstone // reuse of tombstone met before
			}

			return index
		}

		if status == statusOccupied && ht.slots[index] == value {
			return index
		}

		if status == statusDeleted && firstTombstone == -1 {
			firstTombstone = index
		}

		index = (index + ht.step) % ht.size
	}

	// no free slots found but there was a tombstone
	if firstTombstone != -1 {
		return firstTombstone
	}

	return -1
}

func (ht *HashTable) Put(value string) int {
	index := ht.SeekSlot(value)

	if index == -1 {
		return -1
	}

	ht.slots[index] = value
	ht.statuses[index] = statusOccupied

	return index
}

func (ht *HashTable) Find(value string) int {
	index := ht.SeekSlot(value)

	if index != -1 && ht.statuses[index] == statusOccupied && ht.slots[index] == value {
		return index
	}

	return -1
}
