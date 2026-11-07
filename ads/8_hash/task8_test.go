package main

import (
	"testing"
)

func TestHashTable_Init(t *testing.T) {
	t.Run("Valid initialization", func(t *testing.T) {
		ht := Init(HashSize, StepLength)

		if ht.size != HashSize {
			t.Errorf("Expected size %d, got %d", HashSize, ht.size)
		}
		if ht.step != StepLength {
			t.Errorf("Expected step %d, got %d", StepLength, ht.step)
		}
		if len(ht.slots) != HashSize {
			t.Errorf("Expected slots length %d, got %d", HashSize, len(ht.slots))
		}
		if len(ht.statuses) != HashSize {
			t.Errorf("Expected statuses length %d, got %d", HashSize, len(ht.statuses))
		}

		for i, slot := range ht.slots {
			if slot != "" {
				t.Errorf("Expected slot %d to be empty, got %q", i, slot)
			}
			if ht.statuses[i] != statusEmpty {
				t.Errorf("Expected status %d to be statusEmpty, got %d", i, ht.statuses[i])
			}
		}
	})

	t.Run("Initialization with size 1", func(t *testing.T) {
		ht := Init(1, 1)

		if ht.size != 1 || len(ht.slots) != 1 {
			t.Errorf("Expected size and slots length 1, got %d and %d", ht.size, len(ht.slots))
		}
	})
}

func TestHashTable_HashFun(t *testing.T) {
	t.Run("Returns index within bounds", func(t *testing.T) {
		ht := Init(10, 3)

		idx := ht.HashFun("test_value")

		if idx < 0 || idx >= ht.size {
			t.Errorf("HashFun returned out of bounds index: %d", idx)
		}
	})

	t.Run("Consistent hash for same value", func(t *testing.T) {
		ht := Init(10, 3)

		idx1 := ht.HashFun("consistent")
		idx2 := ht.HashFun("consistent")

		if idx1 != idx2 {
			t.Errorf("HashFun returned different indices for the same value: %d and %d", idx1, idx2)
		}
	})
}

func TestHashTable_SeekSlot(t *testing.T) {
	t.Run("Empty table", func(t *testing.T) {
		ht := Init(5, 1)

		idx := ht.SeekSlot("apple")

		if idx < 0 || idx >= 5 {
			t.Errorf("Expected valid index in empty table, got %d", idx)
		}
	})

	t.Run("Full table", func(t *testing.T) {
		ht := Init(3, 1)

		ht.Put("a")
		ht.Put("b")
		ht.Put("c")

		idx := ht.SeekSlot("d")

		if idx != -1 {
			t.Errorf("Expected -1 for full table, got %d", idx)
		}
	})

	t.Run("Finds existing key index even with tombstones along the way", func(t *testing.T) {
		ht := Init(5, 1)

		ht.Put("a")
		idxB := ht.Put("b")
		ht.Put("c")

		// Симулируем удаление "b" строго по логике нашей архитектуры
		ht.slots[idxB] = ""
		ht.statuses[idxB] = statusDeleted

		// SeekSlot для "c" должен пройти через удаленный слот и вернуть точный индекс "c"
		idx := ht.SeekSlot("c")
		if idx == idxB {
			t.Errorf("Expected SeekSlot to find existing 'c', but it returned deleted slot index %d", idxB)
		}
	})
}

func TestHashTable_Put(t *testing.T) {
	t.Run("Put into empty table", func(t *testing.T) {
		ht := Init(5, 1)

		idx := ht.Put("apple")

		if idx < 0 || idx >= 5 {
			t.Errorf("Expected valid index, got %d", idx)
		}
		if ht.slots[idx] != "apple" {
			t.Errorf("Expected slot %d to be 'apple', got %q", idx, ht.slots[idx])
		}
		if ht.statuses[idx] != statusOccupied {
			t.Errorf("Expected status %d to be statusOccupied, got %d", idx, ht.statuses[idx])
		}
	})

	t.Run("Put multiple elements", func(t *testing.T) {
		ht := Init(5, 1)

		ht.Put("apple")
		ht.Put("banana")
		ht.Put("cherry")

		if ht.Find("apple") == -1 || ht.Find("banana") == -1 || ht.Find("cherry") == -1 {
			t.Errorf("Failed to find inserted elements")
		}
	})

	t.Run("Put handles duplicate keys by overwriting same slot", func(t *testing.T) {
		ht := Init(5, 1)

		idx1 := ht.Put("duplicate")
		idx2 := ht.Put("duplicate")

		if idx1 != idx2 {
			t.Errorf("Expected duplicate key to reuse slot %d, but it went to %d", idx1, idx2)
		}
	})

	t.Run("Put into full table", func(t *testing.T) {
		ht := Init(2, 1)

		ht.Put("a")
		ht.Put("b")

		idx := ht.Put("c")

		if idx != -1 {
			t.Errorf("Expected -1 when table is full, got %d", idx)
		}
	})

	t.Run("Put reuses deleted slots (tombstones)", func(t *testing.T) {
		ht := Init(3, 1)

		ht.Put("a")
		idxB := ht.Put("b")
		ht.Put("c")

		// Симулируем удаление элемента вручную через статусы
		ht.slots[idxB] = ""
		ht.statuses[idxB] = statusDeleted

		idxNew := ht.Put("new_element")
		if idxNew != idxB {
			t.Errorf("Expected Put to reuse deleted slot %d, but got %d", idxB, idxNew)
		}
		if ht.slots[idxNew] != "new_element" {
			t.Errorf("Expected slot %d to hold 'new_element', got %q", idxNew, ht.slots[idxNew])
		}
	})
}

func TestHashTable_Find(t *testing.T) {
	t.Run("Find existing element", func(t *testing.T) {
		ht := Init(5, 1)

		ht.Put("apple")

		idx := ht.Find("apple")

		if idx == -1 {
			t.Errorf("Expected to find 'apple', got -1")
		}
		if ht.slots[idx] != "apple" {
			t.Errorf("Expected slot to contain 'apple', got %q", ht.slots[idx])
		}
	})

	t.Run("Find non-existing element in empty table", func(t *testing.T) {
		ht := Init(5, 1)

		idx := ht.Find("apple")

		if idx != -1 {
			t.Errorf("Expected -1 for non-existing element in empty table, got %d", idx)
		}
	})

	t.Run("Find non-existing element with collisions", func(t *testing.T) {
		ht := Init(5, 1)

		ht.Put("a")
		ht.Put("b")
		ht.Put("c")

		idx := ht.Find("z")

		if idx != -1 {
			t.Errorf("Expected -1 for non-existing element, got %d", idx)
		}
	})

	t.Run("Find element correctly passes through tombstones", func(t *testing.T) {
		ht := Init(5, 1)

		idxFirst := ht.Put("first")
		ht.Put("second")

		// Симулируем удаление первого элемента по правилам архитектуры
		ht.slots[idxFirst] = ""
		ht.statuses[idxFirst] = statusDeleted

		// Поиск "second" обязан пройти сквозь надгробие и найти элемент
		if ht.Find("second") == -1 {
			t.Errorf("Expected to find 'second' despite tombstone gap in the table")
		}
	})
}
