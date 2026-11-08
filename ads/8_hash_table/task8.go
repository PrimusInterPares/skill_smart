package main

import (
	// "os"
	"time"
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
	salt     uint32
}

func Init(sz, stp int) HashTable {
	salt := uint32(sz)*uint32(stp)*0x9e3779b9 + uint32(time.Now().UnixNano())

	ht := HashTable{
		size:     sz,
		step:     stp,
		slots:    make([]string, sz),
		statuses: make([]byte, sz),
		salt:     salt,
	}

	return ht
}

func (ht *HashTable) HashFun(value string) int {
	hash := ht.salt

	for i := range len(value) {
		hash ^= uint32(value[i])
		hash *= 16777619 // FNV-1a 32-bit prime multiplier for dynamic bit scattering
	}

	// XORShift mixerto ensure avalanche effect across all bits before taking the modulo
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	// limit hash within HashTable.size and guarantee a positive integer
	return int(hash&0x7FFFFFFF) % ht.size
}

func (ht *HashTable) SeekSlot(value string) int {
	index := ht.HashFun(value)
	firstTombstone := -1

	for range ht.size {
		status := ht.statuses[index]

		if status == statusEmpty && firstTombstone != -1 { // 1. reuse tombstone
			return firstTombstone
		}

		if status == statusEmpty {
			return index // 2. use empty slot
		}

		if status == statusOccupied && ht.slots[index] == value {
			return index // 3. value already exists
		}

		if status == statusDeleted && firstTombstone == -1 {
			firstTombstone = index // store first tombstone met for reuse
		}

		index = (index + ht.step) % ht.size
	}

	return firstTombstone
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


