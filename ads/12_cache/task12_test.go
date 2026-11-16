package main

import (
	"testing"
)

func TestNewNativeCache(t *testing.T) {
	t.Run("returns error on negative size", func(t *testing.T) {
		cache, err := NewNativeCache[int](-1)

		if err == nil {
			t.Fatal("expected error for negative size, got nil")
		}
		if cache != nil {
			t.Fatalf("expected nil cache on error, got %+v", cache)
		}
	})
}

func TestNativeCache_Put(t *testing.T) {
	t.Run("insert into empty cache", func(t *testing.T) {
		cache, err := NewNativeCache[int](5)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("alpha", 10)

		val, ok := cache.Get("alpha")
		if !ok || val != 10 {
			t.Fatalf("expected (10, true), got (%d, %v)", val, ok)
		}
	})

	t.Run("update existing key preserves hits", func(t *testing.T) {
		cache, err := NewNativeCache[int](5)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("alpha", 10)
		cache.Get("alpha")
		cache.Get("alpha")

		cache.Put("alpha", 20)

		val, ok := cache.Get("alpha")
		if !ok || val != 20 {
			t.Fatalf("expected (20, true), got (%d, %v)", val, ok)
		}

		idx := cache.hashFun("alpha")
		if cache.hits[idx] != 3 {
			t.Fatalf("expected 3 hits, got %d", cache.hits[idx])
		}
	})

	t.Run("zero size cache does not panic", func(t *testing.T) {
		cache, err := NewNativeCache[int](0)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("alpha", 10)

		val, ok := cache.Get("alpha")
		if ok || val != 0 {
			t.Fatalf("expected (0, false), got (%d, %v)", val, ok)
		}
	})

	t.Run("reuses deleted slot before empty slot", func(t *testing.T) {
		cache, err := NewNativeCache[int](3)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("a", 1)
		cache.Put("d", 2)
		cache.Put("g", 3)

		cache.Get("d")
		cache.Get("g")

		cache.Put("j", 4)

		val, ok := cache.Get("j")
		if !ok || val != 4 {
			t.Fatalf("expected (4, true), got (%d, %v)", val, ok)
		}

		if cache.slots[1] != "j" {
			t.Fatalf("expected key 'j' at index 1, got %q", cache.slots[1])
		}
	})

	t.Run("reuses actual tombstone before empty slot without duplicating existing key", func(t *testing.T) {
		cache, err := NewNativeCache[int](4)
		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("a", 1)
		cache.Put("e", 2)
		cache.deleteSlot(1)

		cache.Put("e", 20)
		if cache.slots[1] != "" || cache.slotsState[1] != stateDeleted {
			t.Fatalf("expected index 1 to remain deleted when updating 'e', got key=%q state=%v",
				cache.slots[1], cache.slotsState[1])
		}
		if val, ok := cache.Get("e"); !ok || val != 20 {
			t.Fatalf("expected ('e' -> 20, true), got (%d, %v)", val, ok)
		}

		cache.Put("i", 30)
		if cache.slots[1] != "i" || cache.slotsState[1] != stateOccupied {
			t.Fatalf("expected 'i' to reuse tombstone at index 1, got key=%q state=%v", cache.slots[1], cache.slotsState[1])
		}
	})
}

func TestNativeCache_Get(t *testing.T) {
	t.Run("missing key returns zero and false", func(t *testing.T) {
		cache, err := NewNativeCache[string](5)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		val, ok := cache.Get("missing")
		if ok || val != "" {
			t.Fatalf("expected (\"\", false), got (%q, %v)", val, ok)
		}
	})

	t.Run("increments hit counter on every access", func(t *testing.T) {
		cache, err := NewNativeCache[int](5)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("key", 42)

		cache.Get("key")
		cache.Get("key")
		cache.Get("key")

		idx := cache.hashFun("key")
		if cache.hits[idx] != 3 {
			t.Fatalf("expected 3 hits, got %d", cache.hits[idx])
		}
	})

	t.Run("resolves colliding keys properly", func(t *testing.T) {
		cache, err := NewNativeCache[int](3)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("a", 1)
		cache.Put("d", 2)
		cache.Put("g", 3)

		valA, okA := cache.Get("a")
		valD, okD := cache.Get("d")
		valG, okG := cache.Get("g")

		if !okA || valA != 1 {
			t.Fatalf("collision lookup failed for 'a': (%d,%v)", valA, okA)
		}
		if !okD || valD != 2 {
			t.Fatalf("collision lookup failed for 'd': (%d,%v)", valD, okD)
		}
		if !okG || valG != 3 {
			t.Fatalf("collision lookup failed for 'g': (%d,%v)", valG, okG)
		}
	})
}

func TestNativeCache_DeleteSlot(t *testing.T) {
	t.Run("clears slot data and marks as deleted", func(t *testing.T) {
		cache, err := NewNativeCache[string](3)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("a", "value_a")
		cache.Get("a")

		idx := cache.hashFun("a")
		cache.deleteSlot(idx)

		if cache.slotsState[idx] != stateDeleted {
			t.Fatalf("expected stateDeleted, got %v", cache.slotsState[idx])
		}

		if cache.slots[idx] != "" || cache.values[idx] != "" {
			t.Fatalf("slot key or value was not zeroed: key=%q val=%q", cache.slots[idx], cache.values[idx])
		}
		if cache.hits[idx] != 0 {
			t.Fatalf("slot hits was not zeroed: hits=%d", cache.hits[idx])
		}
	})
}

func TestNativeCache_FindMinIdx(t *testing.T) {
	t.Run("finds index with lowest hit count", func(t *testing.T) {
		cache, err := NewNativeCache[int](4)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.hits = []int{5, 2, 8, 4}
		cache.slotsState = []slotState{stateOccupied, stateOccupied, stateOccupied, stateOccupied}

		minIdx := cache.findMinIdx()
		if minIdx != 1 {
			t.Fatalf("expected minIdx 1, got %d", minIdx)
		}
	})

	t.Run("returns first index on tie", func(t *testing.T) {
		cache, err := NewNativeCache[int](4)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.hits = []int{3, 1, 1, 5}
		cache.slotsState = []slotState{stateOccupied, stateOccupied, stateOccupied, stateOccupied}

		minIdx := cache.findMinIdx()
		if minIdx != 1 {
			t.Fatalf("expected minIdx 1, got %d", minIdx)
		}
	})

	t.Run("ignores empty and deleted slots", func(t *testing.T) {
		cache, err := NewNativeCache[int](4)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.hits = []int{0, 5, 0, 3}
		cache.slotsState = []slotState{stateEmpty, stateOccupied, stateDeleted, stateOccupied}

		minIdx := cache.findMinIdx()
		if minIdx != 3 {
			t.Fatalf("expected minIdx 3 (lowest occupied), got %d", minIdx)
		}
	})
}

func TestNativeCache_Eviction(t *testing.T) {
	t.Run("evicts least frequently used key when full", func(t *testing.T) {
		cache, err := NewNativeCache[int](3)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("k1", 100)
		cache.Put("k2", 200)
		cache.Put("k3", 300)

		cache.Get("k1")
		cache.Get("k1")
		cache.Get("k3")

		cache.Put("k4", 400)

		if _, ok := cache.Get("k2"); ok {
			t.Fatal("expected k2 to be evicted")
		}

		val1, ok1 := cache.Get("k1")
		val3, ok3 := cache.Get("k3")
		val4, ok4 := cache.Get("k4")

		if !ok1 || val1 != 100 {
			t.Fatalf("unexpected state for k1: (%d,%v)", val1, ok1)
		}
		if !ok3 || val3 != 300 {
			t.Fatalf("unexpected state for k3: (%d,%v)", val3, ok3)
		}
		if !ok4 || val4 != 400 {
			t.Fatalf("unexpected state for k4: (%d,%v)", val4, ok4)
		}
	})

	t.Run("eviction does not break collision chain", func(t *testing.T) {
		cache, err := NewNativeCache[int](3)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("a", 1)
		cache.Put("d", 2)
		cache.Put("g", 3)

		cache.Get("d")
		cache.Get("g")

		cache.Put("x", 99)

		valD, okD := cache.Get("d")
		valG, okG := cache.Get("g")
		valX, okX := cache.Get("x")

		if !okD || valD != 2 {
			t.Fatalf("broken chain for 'd': (%d,%v)", valD, okD)
		}
		if !okG || valG != 3 {
			t.Fatalf("broken chain for 'g': (%d,%v)", valG, okG)
		}
		if !okX || valX != 99 {
			t.Fatalf("broken chain for 'x': (%d,%v)", valX, okX)
		}
	})

	t.Run("heavy workload with collisions and repeated evictions", func(t *testing.T) {
		cache, err := NewNativeCache[int](4)

		if err != nil {
			t.Fatalf("failed to create cache: %v", err)
		}

		cache.Put("hot_1", 1)
		cache.Put("hot_2", 2)
		cache.Put("cold_1", 3)
		cache.Put("cold_2", 4)

		for range 10 {
			cache.Get("hot_1")
			cache.Get("hot_2")
		}

		cache.Get("cold_1")
		cache.Put("new_1", 5)
		cache.Get("new_1")
		cache.Get("new_1")
		cache.Put("new_2", 6)

		if _, ok := cache.Get("cold_2"); ok {
			t.Fatal("expected cold_2 to be evicted first")
		}

		if _, ok := cache.Get("cold_1"); ok {
			t.Fatal("expected cold_1 to be evicted second")
		}

		for _, key := range []string{"hot_1", "hot_2", "new_1", "new_2"} {
			if _, ok := cache.Get(key); !ok {
				t.Fatalf("expected %q to remain in cache", key)
			}
		}
	})
}
