package main

import (
	"errors"
	"testing"
)

func TestLinkedListDictionary_All(t *testing.T) {
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
