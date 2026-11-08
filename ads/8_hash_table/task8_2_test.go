package main

import (
	"testing"
)

func TestDynHashTable_Init(t *testing.T) {
	t.Run("Valid dynamic initialization", func(t *testing.T) {
		dht := InitDynHashTable(11, 3)

		if dht.count != 0 {
			t.Errorf("Expected initial count 0, got %d", dht.count)
		}
		if dht.size != 11 {
			t.Errorf("Expected initial size 11, got %d", dht.size)
		}
		if len(dht.slots) != 11 {
			t.Errorf("Expected slots length 11, got %d", len(dht.slots))
		}
	})
}

func TestDynHashTable_LoadFactorAndNeedsResize(t *testing.T) {
	t.Run("Load factor changes correctly", func(t *testing.T) {
		dht := InitDynHashTable(10, 1)

		if dht.loadFactor() != 0.0 {
			t.Errorf("Expected load factor 0.0 for empty table, got %f", dht.loadFactor())
		}

		dht.Put("one")
		dht.Put("two")

		expectedLF := 2.0 / 10.0
		if dht.loadFactor() != expectedLF {
			t.Errorf("Expected load factor %f, got %f", expectedLF, dht.loadFactor())
		}
		if dht.needsResize() {
			t.Errorf("Table should not need resize at load factor %f", dht.loadFactor())
		}
	})

	t.Run("Detects when resize is needed", func(t *testing.T) {
		dht := InitDynHashTable(5, 1) // при 4 элементах (4/5 = 0.8) нужен ресайз

		dht.Put("a") // 1/5 = 0.2
		dht.Put("b") // 2/5 = 0.4
		dht.Put("c") // 3/5 = 0.6

		if dht.needsResize() {
			t.Errorf("Table should NOT need resize at 3/5 capacity")
		}

		dht.Put("d") // спровоцирует превышение порога в 0.7

		if !dht.needsResize() {
			t.Errorf("Table SHOULD need resize at 4/5 capacity (LF=%f)", dht.loadFactor())
		}
	})
}

func TestDynHashTable_Resize(t *testing.T) {
	t.Run("Automatic resize to the next prime number", func(t *testing.T) {
		dht := InitDynHashTable(5, 1) // ближайшее простое число после 10 — это 11

		dht.Put("key1") // count=1, size=5
		dht.Put("key2") // count=2, size=5
		dht.Put("key3") // count=3, size=5. LF = 3/5 = 0.6
		dht.Put("key4") // count=4, size=5. LF = 4/5 = 0.8
		dht.Put("key5") // должна вызвать dht.resize() ДО добавления элемента

		if dht.size != 11 {
			t.Errorf("Expected table to resize to next prime number (11), but got size %d", dht.size)
		}
		if dht.count != 5 {
			t.Errorf("Expected count to be 5 after resize and insert, got %d", dht.count)
		}

		// все старые элементы успешно перенеслись и их можно найти на новых индексах
		keys := []string{"key1", "key2", "key3", "key4", "key5"}
		for _, key := range keys {
			if idx := dht.Find(key); idx == -1 {
				t.Errorf("Failed to find key %q after table resize", key)
			}
		}
	})

	t.Run("Duplicate keys do not trigger false count increment or resize", func(t *testing.T) {
		dht := InitDynHashTable(5, 1)

		for range 10 {
			dht.Put("same_key")
		}

		if dht.count != 1 {
			t.Errorf("Expected unique count to be 1 for duplicate inserts, got %d", dht.count)
		}
		if dht.size != 5 {
			t.Errorf("Table should not resize for duplicate keys, expected size 5, got %d", dht.size)
		}
	})

	t.Run("Resize clears out tombstones", func(t *testing.T) {
		dht := InitDynHashTable(5, 1)

		dht.Put("a")
		idxB := dht.Put("b")
		dht.Put("c")

		// симуляция удаления элемента "b"
		dht.slots[idxB] = ""
		dht.statuses[idxB] = statusDeleted
		dht.count-- // уменьшаем счетчик живых элементов

		dht.Put("d")
		dht.Put("e")
		dht.Put("f") // должен произойти resize

		for i, status := range dht.statuses {
			if status == statusDeleted {
				t.Errorf("Found active tombstone at index %d after resize. All tombstones must be cleared!", i)
			}
		}

		if dht.Find("a") == -1 || dht.Find("c") == -1 {
			t.Errorf("Lost valid elements during cleanup-resize")
		}
		if dht.Find("b") != -1 {
			t.Errorf("Found deleted element 'b' after resize")
		}
	})
}

func TestDoubleHashTable_Operations(t *testing.T) {
	t.Run("Basic insertion and search", func(t *testing.T) {
		baseHT := Init(11, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		words := []string{"apple", "banana", "orange", "grape"}
		insertedIndices := make(map[string]int)

		for _, word := range words {
			idx := dht.Put(word)

			if idx == -1 {
				t.Fatalf("failed to insert word %q into an empty table", word)
			}

			insertedIndices[word] = idx
		}

		for _, word := range words {
			idx := dht.Find(word)

			if idx != insertedIndices[word] {
				t.Errorf("for word %q expected index %d, but Find returned %d", word, insertedIndices[word], idx)
			}
		}
	})

	t.Run("Search for non-existent element", func(t *testing.T) {
		baseHT := Init(11, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		if idx := dht.Find("pineapple"); idx != -1 {
			t.Errorf("expected -1 when searching for missing word 'pineapple', but got %d", idx)
		}
	})

	t.Run("Prevent duplicate keys", func(t *testing.T) {
		baseHT := Init(11, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		idx1 := dht.Put("apple")
		idx2 := dht.Put("apple")

		if idx2 != idx1 {
			t.Errorf("on duplicate insert expected old index %d, but got %d", idx1, idx2)
		}
	})

	t.Run("Handling deletion and tombstones", func(t *testing.T) {
		baseHT := Init(11, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		idxApple := dht.Put("apple")
		idxBanana := dht.Put("banana")

		// simulate manual deletion of "apple" to create a tombstone
		dht.statuses[idxApple] = statusDeleted
		dht.slots[idxApple] = ""

		if idx := dht.Find("apple"); idx != -1 {
			t.Errorf("word 'apple' was deleted, but Find still returned index %d", idx)
		}

		// "melon" should safely reuse the tombstone slot
		melonIdx := dht.Put("melon")
		if melonIdx == -1 {
			t.Fatal("failed to insert word 'melon' after creating a tombstone")
		}

		// other elements in the collision chain are still reachable
		if idx := dht.Find("banana"); idx != idxBanana {
			t.Errorf("search for 'banana' broke after tombstone insertion: expected %d, got %d", idxBanana, idx)
		}
	})
}

func TestDoubleHashTable_HashFunctions(t *testing.T) {
	t.Run("Primary hash h1 bounds verification", func(t *testing.T) {
		baseHT := Init(19, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		h1Val := dht.h1("test_string")
		if h1Val < 0 || h1Val >= dht.size {
			t.Errorf("h1 returned index out of bounds: %d for table size %d", h1Val, dht.size)
		}
	})

	t.Run("Secondary hash h2 step size validation", func(t *testing.T) {
		baseHT := Init(19, StepLength)
		dht := DoubleHashTable{HashTable: &baseHT}

		testStrings := []string{"a", "b", "hello", "very_long_string_for_hash_testing", ""}
		for _, testStr := range testStrings {
			step := dht.h2(testStr)

			if step <= 0 {
				t.Errorf("h2 returned invalid step size %d (must be > 0) for string %q", step, testStr)
			}
			if step >= dht.size {
				t.Errorf("h2 returned step size %d which is >= table size %d for string %q", step, dht.size, testStr)
			}
		}
	})
}
