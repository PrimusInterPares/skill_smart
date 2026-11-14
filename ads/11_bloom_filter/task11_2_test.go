package main

import (
	"errors"
	"testing"
)

func TestMergeBloomFilters(t *testing.T) {
	t.Run("error when no filters provided", func(t *testing.T) {
		got, err := MergeBloomFilters()

		if !errors.Is(err, ErrNoFiltersProvided) {
			t.Errorf("got error %v; want %v", err, ErrNoFiltersProvided)
		}
		if got != nil {
			t.Errorf("expected nil result on error")
		}
	})

	t.Run("error when nil filter is passed first or in the middle", func(t *testing.T) {
		valid := NewBloomFilter(32)

		if _, err := MergeBloomFilters(nil, valid); !errors.Is(err, ErrNilFilter) {
			t.Errorf("nil as first arg: got %v; want %v", err, ErrNilFilter)
		}

		if _, err := MergeBloomFilters(valid, nil); !errors.Is(err, ErrNilFilter) {
			t.Errorf("nil as second arg: got %v; want %v", err, ErrNilFilter)
		}
	})

	t.Run("error when filter sizes differ", func(t *testing.T) {
		bf32 := NewBloomFilter(32)
		bf64 := NewBloomFilter(64)

		if _, err := MergeBloomFilters(bf32, bf64); !errors.Is(err, ErrFilterSizesDiffer) {
			t.Errorf("got error %v; want %v", err, ErrFilterSizesDiffer)
		}
	})

	t.Run("error when filter size is zero or negative", func(t *testing.T) {
		zeroBf := &BloomFilter{filter_len: 0}
		negBf := &BloomFilter{filter_len: -10}

		if _, err := MergeBloomFilters(zeroBf); !errors.Is(err, ErrInvalidFilterSize) {
			t.Errorf("zero size: got %v; want %v", err, ErrInvalidFilterSize)
		}

		if _, err := MergeBloomFilters(negBf); !errors.Is(err, ErrInvalidFilterSize) {
			t.Errorf("negative size: got %v; want %v", err, ErrInvalidFilterSize)
		}
	})

	t.Run("error when bits slice length does not match filter_len", func(t *testing.T) {
		corrupted := &BloomFilter{
			filter_len: 64,
			bits:       make([]uint32, 1),
		}
		valid := NewBloomFilter(64)

		if _, err := MergeBloomFilters(valid, corrupted); !errors.Is(err, ErrIrrelevantLength) {
			t.Errorf("got error %v; want %v", err, ErrIrrelevantLength)
		}
	})

	t.Run("merge single filter returns independent copy", func(t *testing.T) {
		bf := NewBloomFilter(32)
		bf.Add("alpha")

		merged, err := MergeBloomFilters(bf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !merged.IsValue("alpha") {
			t.Errorf("merged filter missing element from source")
		}

		merged.Add("beta")
		if bf.IsValue("beta") {
			t.Errorf("modifying merged filter mutated original source filter")
		}
	})

	t.Run("merge multiple populated filters across multi-word boundaries", func(t *testing.T) {
		bf1 := NewBloomFilter(100)
		bf2 := NewBloomFilter(100)
		bf3 := NewBloomFilter(100)

		bf1.Add("node1_key")
		bf2.Add("node2_key")
		bf3.Add("node3_key")

		merged, err := MergeBloomFilters(bf1, bf2, bf3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, key := range []string{"node1_key", "node2_key", "node3_key"} {
			if !merged.IsValue(key) {
				t.Errorf("merged filter missing key %q", key)
			}
		}

		if bf1.IsValue("node2_key") || bf2.IsValue("node3_key") {
			t.Errorf("source filters were mutated during merge")
		}
	})

	t.Run("merge handles empty filters with uninitialized nil bits slice", func(t *testing.T) {
		emptyDirect := &BloomFilter{filter_len: 64}
		populated := NewBloomFilter(64)
		populated.Add("present_key")

		merged, err := MergeBloomFilters(emptyDirect, populated)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !merged.IsValue("present_key") {
			t.Errorf("merged filter lost key when merged with nil-bits filter")
		}
	})
}
