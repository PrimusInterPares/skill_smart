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

func TestCountingBloomFilter_NewCountingBloomFilter(t *testing.T) {
	t.Run("creates filter with default size when provided size is zero or negative", func(t *testing.T) {
		tests := []struct {
			size int
			want int
		}{
			{0, 1},
			{-10, 1},
			{10, 10},
		}

		for _, test := range tests {
			cbf := NewCountingBloomFilter(test.size)

			if cbf.filterLen != test.want {
				t.Errorf("NewCountingBloomFilter(%d) = %d; want %d", test.size, cbf.filterLen, test.want)
			}
			if len(cbf.counters) != countersLen(test.want) {
				t.Errorf("counters length for size %d is incorrect: got %d, want %d", test.want, len(cbf.counters), countersLen(test.want))
			}
		}
	})
}

func TestCountingBloomFilter_GetCount(t *testing.T) {
	t.Run("returns zero for out-of-bounds index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		if got := cbf.getCount(15); got != 0 {
			t.Errorf("getCount(15) = %d; want 0", got)
		}
	})

	t.Run("returns zero for uninitialized counters", func(t *testing.T) {
		cbf := &CountingBloomFilter{filterLen: 10, counters: nil}

		if got := cbf.getCount(5); got != 0 {
			t.Errorf("getCount(5) = %d; want 0", got)
		}
	})

	t.Run("returns correct count for valid index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(5, 7)

		if got := cbf.getCount(5); got != 7 {
			t.Errorf("getCount(5) = %d; want 7", got)
		}
	})
}

func TestCountingBloomFilter_SetCount(t *testing.T) {
	t.Run("sets count for valid index without corrupting adjacent nibbles", func(t *testing.T) {
		cbf := NewCountingBloomFilter(16)
		boundaryIndices := []uint{0, 7, 8, 15}

		for _, idx := range boundaryIndices {
			cbf.setCount(idx, MaxCounterVal)
		}

		for _, idx := range boundaryIndices {
			if got := cbf.getCount(idx); got != MaxCounterVal {
				t.Errorf("getCount(%d) = %d; want %d", idx, got, MaxCounterVal)
			}
		}

		for _, untouched := range []uint{1, 6, 9, 14} {
			if got := cbf.getCount(untouched); got != 0 {
				t.Errorf("adjacent counter at %d was corrupted: got %d; want 0", untouched, got)
			}
		}
	})

	t.Run("does nothing for out-of-bounds index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(15, 7)

		if got := cbf.getCount(15); got != 0 {
			t.Errorf("getCount(15) = %d; want 0", got)
		}
	})

	t.Run("initializes counters if they are nil", func(t *testing.T) {
		cbf := &CountingBloomFilter{filterLen: 10, counters: nil}

		cbf.setCount(5, 7)

		if len(cbf.counters) != countersLen(10) {
			t.Errorf("counters length is incorrect: got %d, want %d", len(cbf.counters), countersLen(10))
		}
		if got := cbf.getCount(5); got != 7 {
			t.Errorf("getCount(5) = %d; want 7", got)
		}
	})
}

func TestCountingBloomFilter_IncCount(t *testing.T) {
	t.Run("increments count for valid index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(5, 7)
		cbf.incCount(5)

		if got := cbf.getCount(5); got != 8 {
			t.Errorf("getCount(5) = %d; want 8", got)
		}
	})

	t.Run("does nothing for out-of-bounds index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.incCount(15)

		if got := cbf.getCount(15); got != 0 {
			t.Errorf("getCount(15) = %d; want 0", got)
		}
	})

	t.Run("saturates at max value without overflow", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(5, MaxCounterVal)
		cbf.incCount(5)

		if got := cbf.getCount(5); got != MaxCounterVal {
			t.Errorf("getCount(5) = %d; want %d", got, MaxCounterVal)
		}
	})
}

func TestCountingBloomFilter_DecCount(t *testing.T) {
	t.Run("decrements count for valid index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(5, 7)
		cbf.decCount(5)

		if got := cbf.getCount(5); got != 6 {
			t.Errorf("getCount(5) = %d; want 6", got)
		}
	})

	t.Run("does nothing for out-of-bounds index", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.decCount(15)

		if got := cbf.getCount(15); got != 0 {
			t.Errorf("getCount(15) = %d; want 0", got)
		}
	})

	t.Run("does nothing if count is already zero", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.decCount(5)

		if got := cbf.getCount(5); got != 0 {
			t.Errorf("getCount(5) = %d; want 0", got)
		}
	})

	t.Run("does nothing if count is saturated at max value", func(t *testing.T) {
		cbf := NewCountingBloomFilter(10)

		cbf.setCount(5, MaxCounterVal)
		cbf.decCount(5)

		if got := cbf.getCount(5); got != MaxCounterVal {
			t.Errorf("getCount(5) = %d; want %d", got, MaxCounterVal)
		}
	})
}

func TestCountingBloomFilter_PublicAPI(t *testing.T) {
	t.Run("add, check and remove element", func(t *testing.T) {
		cbf := NewCountingBloomFilter(32)

		if cbf.IsValue("alpha") {
			t.Errorf("expected false before adding element")
		}

		cbf.Add("alpha")
		if !cbf.IsValue("alpha") {
			t.Errorf("expected true after adding element")
		}

		cbf.Remove("alpha")
		if cbf.IsValue("alpha") {
			t.Errorf("expected false after removing element")
		}
	})

	t.Run("removing one element does not remove colliding element", func(t *testing.T) {
		cbf := NewCountingBloomFilter(32)
		items := []string{
			"0123456789", "1234567890", "2345678901", "3456789012", "4567890123",
			"5678901234", "6789012345", "7890123456", "8901234567", "9012345678",
		}

		for _, item := range items {
			cbf.Add(item)
		}

		cbf.Remove(items[0])

		for _, remaining := range items[1:] {
			if !cbf.IsValue(remaining) {
				t.Errorf("false negative for %q after removing another element", remaining)
			}
		}
	})

	t.Run("removing duplicate additions requires matching number of removals", func(t *testing.T) {
		cbf := NewCountingBloomFilter(32)

		cbf.Add("dup")
		cbf.Add("dup")

		cbf.Remove("dup")
		if !cbf.IsValue("dup") {
			t.Errorf("expected element to remain after first removal")
		}

		cbf.Remove("dup")
		if cbf.IsValue("dup") {
			t.Errorf("expected element to be absent after second removal")
		}
	})

	t.Run("removing guaranteed absent element does not corrupt existing counters", func(t *testing.T) {
		cbf := NewCountingBloomFilter(64)
		cbf.Add("existing")

		cbf.Remove("non_existent_key")

		if !cbf.IsValue("existing") {
			t.Errorf("removing non-existent key corrupted existing element")
		}
	})

	t.Run("saturated counter prevents false negatives after removal", func(t *testing.T) {
		cbf := NewCountingBloomFilter(32)

		for range MaxCounterVal + 5 {
			cbf.Add("heavy_key")
		}

		cbf.Remove("heavy_key")

		if !cbf.IsValue("heavy_key") {
			t.Errorf("saturated element must remain present to avoid false negatives")
		}
	})
}
