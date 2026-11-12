package main

import (
	"testing"
)

func TestPowerSet_CartesianProduct(t *testing.T) {
	t.Run("non-empty sets", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("A")
		set1.Put("B")
		set2.Put("1")
		set2.Put("2")

		prod := set1.CartesianProduct(set2)
		if prod.Size() != 4 {
			t.Errorf("Expected cartesian product size 4, got %d", prod.Size())
		}
		if !prod.Get("(A, 1)") || !prod.Get("(A, 2)") {
			t.Errorf("Expected cartesian product to contain pairs for 'A'")
		}
		if !prod.Get("(B, 1)") || !prod.Get("(B, 2)") {
			t.Errorf("Expected cartesian product to contain pairs for 'B'")
		}
	})

	t.Run("one empty set", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("A")
		set1.Put("B")

		prod := set1.CartesianProduct(set2)
		if prod.Size() != 0 {
			t.Errorf("Expected cartesian product with an empty set to be empty, got size %d", prod.Size())
		}
	})

	t.Run("both empty sets", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		prod := set1.CartesianProduct(set2)
		if prod.Size() != 0 {
			t.Errorf("Expected cartesian product of two empty sets to be empty, got size %d", prod.Size())
		}
	})
}
