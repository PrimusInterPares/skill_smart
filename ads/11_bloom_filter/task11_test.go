package main

import (
	"testing"
)

func TestBloomFilter_Hashes(t *testing.T) {
	bf := NewBloomFilter(32)

	t.Run("hash values are within bounds [0, filter_len)", func(t *testing.T) {

		inputs := []string{"", "a", "0123456789", "hello world", "special!@#$%^&*()"}
		for _, s := range inputs {
			h1 := bf.Hash1(s)
			h2 := bf.Hash2(s)

			if h1 < 0 || h1 >= bf.filter_len {
				t.Errorf("Hash1(%q) = %d; want within [0, %d)", s, h1, bf.filter_len)
			}
			if h2 < 0 || h2 >= bf.filter_len {
				t.Errorf("Hash2(%q) = %d; want within [0, %d)", s, h2, bf.filter_len)
			}
		}
	})

	t.Run("hashes are deterministic", func(t *testing.T) {
		s := "0123456789"

		h1a, h1b := bf.Hash1(s), bf.Hash1(s)
		if h1a != h1b {
			t.Errorf("Hash1 is not deterministic for %q", s)
		}

		h2a, h2b := bf.Hash2(s), bf.Hash2(s)
		if h2a != h2b {
			t.Errorf("Hash2 is not deterministic for %q", s)
		}
	})

	t.Run("empty string produces zero hash", func(t *testing.T) {
		if got := bf.Hash1(""); got != 0 {
			t.Errorf("Hash1(\"\") = %d; want 0", got)
		}
		if got := bf.Hash2(""); got != 0 {
			t.Errorf("Hash2(\"\") = %d; want 0", got)
		}
	})

	t.Run("zero or negative filter_len safety", func(t *testing.T) {
		zeroBf := &BloomFilter{filter_len: 0}

		if got := zeroBf.Hash1("test"); got != 0 {
			t.Errorf("Hash1 with filter_len=0 got %d; want 0", got)
		}

		negBf := &BloomFilter{filter_len: -10}

		if got := negBf.Hash2("test"); got != 0 {
			t.Errorf("Hash2 with filter_len=-10 got %d; want 0", got)
		}
	})
}

func TestBloomFilter_AddAndIsValue(t *testing.T) {
	t.Run("empty filter returns false", func(t *testing.T) {
		bf := NewBloomFilter(32)

		if bf.IsValue("0123456789") {
			t.Errorf("IsValue on empty filter returned true")
		}
		if bf.IsValue("") {
			t.Errorf("IsValue(\"\") on empty filter returned true")
		}
	})

	t.Run("task specification 10 cyclic strings no false negatives", func(t *testing.T) {
		bf := NewBloomFilter(32)
		items := []string{
			"0123456789",
			"1234567890",
			"2345678901",
			"3456789012",
			"4567890123",
			"5678901234",
			"6789012345",
			"7890123456",
			"8901234567",
			"9012345678",
		}

		for _, item := range items {
			bf.Add(item)
		}

		for _, item := range items {
			if !bf.IsValue(item) {
				t.Errorf("false negative for added item %q", item)
			}
		}
	})

	t.Run("idempotency of multiple additions", func(t *testing.T) {
		bf := NewBloomFilter(32)
		item := "duplicate_key"

		bf.Add(item)
		bitsAfterFirst := append([]uint32(nil), bf.bits...)

		bf.Add(item)
		bf.Add(item)

		if len(bf.bits) != len(bitsAfterFirst) {
			t.Fatalf("bits slice length changed")
		}

		for i := range bf.bits {
			if bf.bits[i] != bitsAfterFirst[i] {
				t.Errorf("bits[%d] changed after duplicate Add: got %b, want %b", i, bf.bits[i], bitsAfterFirst[i])
			}
		}
	})

	t.Run("empty string edge case", func(t *testing.T) {
		bf := NewBloomFilter(32)

		bf.Add("")

		if !bf.IsValue("") {
			t.Errorf("expected IsValue(\"\") to be true after Add(\"\")")
		}
	})

	t.Run("tricky multi-word boundaries (>32 bits)", func(t *testing.T) {
		bf := NewBloomFilter(100)
		items := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

		for _, item := range items {
			bf.Add(item)
		}

		for _, item := range items {
			if !bf.IsValue(item) {
				t.Errorf("false negative on 100-bit filter for %q", item)
			}
		}

		if bf.IsValue("non_existent_omega") {
			t.Errorf("unexpected collision on large 100-bit filter")
		}
	})

	t.Run("filter saturation causes expected false positives", func(t *testing.T) {
		bf := NewBloomFilter(1)
		bf.Add("anything")

		if !bf.IsValue("completely_different_string") {
			t.Errorf("expected false positive when 1-bit filter is saturated")
		}
	})

	t.Run("direct struct literal initialization without NewBloomFilter", func(t *testing.T) {
		bf := &BloomFilter{filter_len: 32}

		if bf.IsValue("test") {
			t.Errorf("expected false on uninitialized bits slice")
		}
	})
}

func TestBloomFilter_BitOperations(t *testing.T) {
	t.Run("boundary bit indices 0, 31, 32, 63", func(t *testing.T) {
		bf := NewBloomFilter(64)
		indices := []uint{0, 31, 32, 63}

		for _, idx := range indices {
			if bf.testBit(idx) {
				t.Errorf("bit %d should initially be 0", idx)
			}

			bf.setBit(idx)

			if !bf.testBit(idx) {
				t.Errorf("bit %d should be 1 after setBit", idx)
			}
		}

		if bf.testBit(1) || bf.testBit(30) {
			t.Errorf("adjacent bits were corrupted")
		}
		if bf.testBit(33) || bf.testBit(62) {
			t.Errorf("adjacent bits were corrupted")
		}
	})

	t.Run("out of bounds bit indices do not panic", func(t *testing.T) {
		bf := NewBloomFilter(32)

		bf.setBit(32)
		bf.setBit(1000)

		if bf.testBit(32) || bf.testBit(1000) {
			t.Errorf("out of bounds testBit must return false")
		}
	})
}
