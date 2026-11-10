package main

import (
	"testing"
)

func TestNativeDictionary_Put(t *testing.T) {
	t.Run("New key insertion", func(t *testing.T) {
		nd := Init[int](10)

		nd.Put("key1", 42)

		v, err := nd.Get("key1")
		if err != nil {
			t.Errorf("Expected key1 to be found, got error: %v", err)
		}
		if v != 42 {
			t.Errorf("Expected value to be 42, got %d", v)
		}
	})

	t.Run("Existing key update", func(t *testing.T) {
		nd := Init[int](10)

		nd.Put("key1", 10)
		nd.Put("key1", 99) // update value

		v, err := nd.Get("key1")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 99 {
			t.Errorf("Expected value to be updated to 99, got %d", v)
		}
	})
}

func TestNativeDictionary_IsKey(t *testing.T) {
	t.Run("Present key", func(t *testing.T) {
		nd := Init[int](10)

		nd.Put("present", 100)

		if !nd.IsKey("present") {
			t.Error("Expected IsKey to return true for present key")
		}
	})

	t.Run("Absent key", func(t *testing.T) {
		nd := Init[int](10)

		if nd.IsKey("absent") {
			t.Error("Expected IsKey to return false for absent key")
		}
	})
}

func TestNativeDictionary_Get(t *testing.T) {
	t.Run("Existing key extraction", func(t *testing.T) {
		nd := Init[int](10)

		nd.Put("key1", 42)

		v, err := nd.Get("key1")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 42 {
			t.Errorf("Expected value 42, got %d", v)
		}
	})

	t.Run("Absent key extraction", func(t *testing.T) {
		nd := Init[int](10)

		v, err := nd.Get("missing")
		if err == nil {
			t.Error("Expected error for absent key, got nil")
		}
		if v != 0 { // zero value for int
			t.Errorf("Expected default value 0, got %d", v)
		}
	})
}

func TestNativeDictionary_MixedOperations(t *testing.T) {
	t.Run("Put and Get cycle", func(t *testing.T) {
		nd := Init[string](20)

		nd.Put("new_key", "value1")

		v, err := nd.Get("new_key")
		if err != nil || v != "value1" {
			t.Errorf("Expected 'value1', got '%s' (err: %v)", v, err)
		}

		nd.Put("new_key", "value2")

		v, err = nd.Get("new_key")
		if err != nil || v != "value2" {
			t.Errorf("Expected updated 'value2', got '%s' (err: %v)", v, err)
		}
	})

	t.Run("Full capacity handling", func(t *testing.T) {
		sz := 5
		nd := Init[int](sz)

		for i := range sz {
			nd.Put("k"+string(rune('0'+i)), i*10)
		}

		for i := range sz {
			v, err := nd.Get("k" + string(rune('0'+i)))
			if err != nil {
				t.Errorf("Unexpected error for key %d: %v", i, err)
			}
			if v != i*10 {
				t.Errorf("Expected %d, got %d", i*10, v)
			}
		}
	})
}
