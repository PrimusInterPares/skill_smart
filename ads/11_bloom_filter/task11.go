package main

import (
// "os"
)

const (
	hash1MagicMultiplier = 17
	hash2MagicMultiplier = 223
)

type BloomFilter struct {
	filter_len int
	bits       []uint32
}

func NewBloomFilter(flen int) *BloomFilter {
	filterLen := flen
	if filterLen <= 0 {
		filterLen = 1
	}

	return &BloomFilter{
		filter_len: filterLen,
		bits:       initBits((filterLen + 31) >> 5), // divide into sections of 32 bits each
	}
}

func initBits(size int) []uint32 {
	return make([]uint32, size)
}

func (bf *BloomFilter) setBit(bitIdx uint) {
	if bitIdx >= uint(bf.filter_len) {
		return
	}

	if len(bf.bits) == 0 && bf.filter_len > 0 {
		bf.bits = initBits((bf.filter_len + 31) >> 5)
	}

	bf.bits[bitIdx>>5] |= (uint32(1) << (bitIdx & 31)) // bitIdx % 32
}

func (bf *BloomFilter) testBit(bitIdx uint) bool {
	if bitIdx >= uint(bf.filter_len) {
		return false
	}

	return (bf.bits[bitIdx>>5] & (uint32(1) << (bitIdx & 31))) != 0
}

func (bf *BloomFilter) commonHash(s string, magicMultiplier int) int {
	if bf.filter_len <= 0 {
		return 0
	}

	sum := 0
	for _, char := range s {
		code := int(char)
		sum = (sum*magicMultiplier + code) % bf.filter_len
	}

	return sum

}

func (bf *BloomFilter) Hash1(s string) int {
	return bf.commonHash(s, hash1MagicMultiplier)
}

func (bf *BloomFilter) Hash2(s string) int {
	return bf.commonHash(s, hash2MagicMultiplier)
}

func (bf *BloomFilter) Add(s string) {
	h1, h2 := bf.Hash1(s), bf.Hash2(s)

	bf.setBit(uint(h1))
	bf.setBit(uint(h2))
}

func (bf *BloomFilter) IsValue(s string) bool {
	if len(bf.bits) == 0 && bf.filter_len > 0 {
		return false
	}

	h1, h2 := bf.Hash1(s), bf.Hash2(s)

	return bf.testBit(uint(h1)) && bf.testBit(uint(h2))
}
