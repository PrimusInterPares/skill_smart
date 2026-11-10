package main

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var ErrKeyNotFound = errors.New("key is not found")

type NativeDictionary[T any] struct {
	size   int
	slots  []string
	values []T
	salt   uint32
}

func Init[T any](sz int) NativeDictionary[T] {
	nd := NativeDictionary[T]{
		size:   sz,
		slots:  nil,
		values: nil,
		salt:   uint32(sz)*0x9e3779b9 + uint32(time.Now().UnixNano()),
	}
	nd.slots = make([]string, sz)
	nd.values = make([]T, sz)

	return nd
}

func (nd *NativeDictionary[T]) HashFun(value string) int {
	hash := nd.salt

	for i := range len(value) {
		hash ^= uint32(value[i])
		hash *= 16777619
	}

	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return int(hash&0x7FFFFFFF) % nd.size
}

func (nd *NativeDictionary[T]) IsKey(key string) bool {
	startIndex := nd.HashFun(key)
	index := startIndex

	for nd.slots[index] != "" {
		if nd.slots[index] == key {
			return true
		}

		index = (index + 1) % nd.size

		if index == startIndex {
			return false
		}
	}

	return false
}

func (nd *NativeDictionary[T]) Get(key string) (T, error) {
	var zero T

	startIndex := nd.HashFun(key)
	index := startIndex

	for nd.slots[index] != "" {
		if nd.slots[index] == key {
			return nd.values[index], nil
		}

		index = (index + 1) % nd.size

		if index == startIndex {
			return zero, ErrKeyNotFound
		}
	}

	return zero, ErrKeyNotFound
}

func (nd *NativeDictionary[T]) Put(key string, value T) {
	startIndex := nd.HashFun(key)
	index := startIndex

	for nd.slots[index] != key {
		if nd.slots[index] == "" {
			nd.slots[index] = key

			break
		}

		index = (index + 1) % nd.size

		if index == startIndex {
			return
		}
	}

	nd.values[index] = value
}


