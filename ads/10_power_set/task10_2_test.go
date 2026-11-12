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

func TestPowerSet_MultiIntersection(t *testing.T) {
	t.Run("three sets with common elements", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()
		set3 := Init[string]()

		set1.Put("A")
		set1.Put("B")
		set1.Put("C")

		set2.Put("B")
		set2.Put("C")
		set2.Put("D")

		set3.Put("C")
		set3.Put("E")

		sets := []PowerSet[string]{set1, set2, set3}
		res := MultiIntersection(sets)

		if res.Size() != 1 || !res.Get("C") {
			t.Errorf("Expected intersection to contain only 'C', got size %d", res.Size())
		}
	})

	t.Run("three sets with no common elements", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()
		set3 := Init[string]()

		set1.Put("A")
		set1.Put("B")

		set2.Put("B")
		set2.Put("C")

		set3.Put("C")
		set3.Put("D")

		sets := []PowerSet[string]{set1, set2, set3}
		res := MultiIntersection(sets)

		if res.Size() != 0 {
			t.Errorf("Expected empty intersection, got size %d", res.Size())
		}
	})

	t.Run("empty slice of sets", func(t *testing.T) {
		var sets []PowerSet[string]
		res := MultiIntersection(sets)

		if res.Size() != 0 {
			t.Errorf("Expected empty set for empty input slice, got size %d", res.Size())
		}
	})

	t.Run("single set in slice", func(t *testing.T) {
		set1 := Init[string]()
		set1.Put("X")
		set1.Put("Y")

		sets := []PowerSet[string]{set1}
		res := MultiIntersection(sets)

		if res.Size() != 2 || !res.Get("X") {
			t.Errorf("Expected copy of the single set with size 2, got size %d", res.Size())
		}
		if !res.Get("Y") {
			t.Errorf("Expected 'Y' to be present in the resulting copy")
		}
	})
}
