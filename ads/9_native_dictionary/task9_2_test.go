package main

import (
	"errors"
	"testing"
)

func TestLinkedListDictionary_IsKey(t *testing.T) {
	t.Run("Empty dictionary behavior", func(t *testing.T) {
		lld := NewLinkedListDictionary[string, int]()

		if lld.IsKey("apple") {
			t.Error("Expected IsKey to return false for empty dict")
		}

		_, err := lld.Get("apple")
		if !errors.Is(err, ErrEmptyDictionary) {
			t.Errorf("Expected ErrEmptyDictionary, got %v", err)
		}
	})
}

func TestLinkedListDictionary_Put(t *testing.T) {
	t.Run("Put insertions and sorting order", func(t *testing.T) {
		lld := NewLinkedListDictionary[string, int]()

		lld.Put("cherry", 30)
		lld.Put("apple", 10)
		lld.Put("banana", 20)

		val, err := lld.Get("banana")
		if err != nil || val != 20 {
			t.Errorf("Expected 20, got %d (err: %v)", val, err)
		}

		if lld.head == nil || lld.head.key != "apple" {
			t.Errorf("Expected head to be 'apple', got %v", lld.head)
		}
		if lld.head.next == nil || lld.head.next.key != "banana" {
			t.Errorf("Expected second node to be 'banana'")
		}
		if lld.head.next.next == nil || lld.head.next.next.key != "cherry" {
			t.Errorf("Expected third node to be 'cherry'")
		}
		if lld.count != 3 {
			t.Errorf("Expected count to be 3, got %d", lld.count)
		}
	})

	t.Run("Early stop optimization on missing keys", func(t *testing.T) {
		lld := NewLinkedListDictionary[string, int]()

		lld.Put("apple", 10)
		lld.Put("cherry", 30)

		if lld.IsKey("banana") {
			t.Error("Expected IsKey to return false for 'banana'")
		}

		_, err := lld.Get("banana")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected ErrKeyNotFound for missing middle key, got %v", err)
		}
	})

	t.Run("Put updates existing key value", func(t *testing.T) {
		lld := NewLinkedListDictionary[string, int]()

		lld.Put("banana", 20)
		lld.Put("banana", 99)

		val, _ := lld.Get("banana")
		if val != 99 {
			t.Errorf("Expected updated value 99, got %d", val)
		}
		if lld.count != 1 {
			t.Errorf("Count should remain 1 after update, got %d", lld.count)
		}

		lld.Put("cherry", 30) // banana (head) -> cherry
		lld.Put("banana", 100)

		val, _ = lld.Get("banana")
		if val != 100 {
			t.Errorf("Expected updated head value 100, got %d", val)
		}
		if lld.count != 2 {
			t.Errorf("Count should remain 2 after updating head, got %d", lld.count)
		}

		if lld.head.next == nil || lld.head.next.key != "cherry" {
			t.Error("List structure was broken after head update")
		}
	})
}

func TestLinkedListDictionary_Delete(t *testing.T) {
	t.Run("Delete operations", func(t *testing.T) {
		lld := NewLinkedListDictionary[string, int]()

		lld.Put("apple", 10)
		lld.Put("banana", 20)
		lld.Put("cherry", 30)

		err := lld.Delete("banana")
		if err != nil {
			t.Errorf("Unexpected error on delete: %v", err)
		}
		if lld.IsKey("banana") {
			t.Error("Key 'banana' should have been deleted")
		}
		if lld.count != 2 {
			t.Errorf("Expected count 2, got %d", lld.count)
		}

		err = lld.Delete("apple")
		if err != nil {
			t.Errorf("Unexpected error on deleting head: %v", err)
		}
		if lld.head.key != "cherry" {
			t.Errorf("Expected new head to be 'cherry', got '%s'", lld.head.key)
		}

		_ = lld.Delete("cherry")
		if lld.count != 0 || lld.head != nil {
			t.Error("Dictionary should be completely empty now")
		}

		err = lld.Delete("missing")
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected ErrKeyNotFound on deleting non-existent key, got %v", err)
		}
	})
}

// ======================================================================================================

func TestBinaryTrieDictionary_IsKey(t *testing.T) {
	t.Run("Key is missing in empty dictionary", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		if btd.IsKey(7) {
			t.Error("Expected IsKey to return false for empty dictionary")
		}
	})

	t.Run("Key is present", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(7, "seven")

		if !btd.IsKey(7) {
			t.Error("Expected IsKey to return true for present key")
		}
	})

	t.Run("Key is missing but shares common prefix branch", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(5, "apple")

		if btd.IsKey(4) {
			t.Error("Expected IsKey to return false for key 4 (only prefix exists)")
		}
	})
}

func TestBinaryTrieDictionary_Get(t *testing.T) {
	t.Run("Get value of existing key", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(5, "apple")

		val, err := btd.Get(5)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if val != "apple" {
			t.Errorf("Expected 'apple', got '%s'", val)
		}
	})

	t.Run("Get error for absent key", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		_, err := btd.Get(999)
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected ErrKeyNotFound, got %v", err)
		}
	})
}

func TestBinaryTrieDictionary_Put(t *testing.T) {
	t.Run("Insert new key", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(42, "value42")

		if btd.count != 1 {
			t.Errorf("Expected count to be 1, got %d", btd.count)
		}
	})

	t.Run("Update existing key value", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[int](32)

		btd.Put(100, 10)
		btd.Put(100, 99)

		if btd.count != 1 {
			t.Errorf("Count should remain 1 after updating the same key, got %d", btd.count)
		}
	})
}

func TestBinaryTrieDictionary_Delete(t *testing.T) {
	t.Run("Delete existing key", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(5, "apple")

		err := btd.Delete(5)
		if err != nil {
			t.Errorf("Unexpected error on delete: %v", err)
		}
		if btd.count != 0 {
			t.Errorf("Expected count to be 0, got %d", btd.count)
		}
	})

	t.Run("Delete absent key", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		err := btd.Delete(999)
		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected ErrKeyNotFound on deleting non-existent key, got %v", err)
		}
	})
}

func TestBinaryTrieDictionary_Integration(t *testing.T) {
	t.Run("Put, Get and Delete full cycle with neighbor preservation", func(t *testing.T) {
		btd := NewBinaryTrieDictionary[string](32)

		btd.Put(5, "apple")
		btd.Put(4, "banana")

		btd.Put(5, "new_apple")
		val, _ := btd.Get(5)

		if val != "new_apple" {
			t.Errorf("Expected updated value 'new_apple', got '%s'", val)
		}

		_ = btd.Delete(5)
		if btd.IsKey(5) {
			t.Error("Key 5 should be reported as missing after deletion")
		}

		_, err := btd.Get(5)
		if !errors.Is(err, ErrKeyNotFound) {
			t.Error("Get for key 5 should return ErrKeyNotFound")
		}
		if !btd.IsKey(4) {
			t.Error("Neighbor key 4 should still be present in the tree")
		}

		neighborVal, _ := btd.Get(4)
		if neighborVal != "banana" {
			t.Errorf("Expected neighbor value to remain 'banana', got '%s'", neighborVal)
		}
	})
}


