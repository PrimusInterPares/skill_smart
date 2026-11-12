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

func TestBag_Put(t *testing.T) {
	t.Run("unique elements", func(t *testing.T) {
		bag := InitBag[string](10)

		bag.Put("apple")
		bag.Put("banana")

		if bag.Size() != 2 {
			t.Errorf("Expected unique size 2, got %d", bag.Size())
		}
	})

	t.Run("duplicate elements", func(t *testing.T) {
		bag := InitBag[string](10)

		bag.Put("apple")
		bag.Put("apple")
		bag.Put("apple")

		if bag.Size() != 1 {
			t.Errorf("Expected unique size 1 for duplicates, got %d", bag.Size())
		}

		elements := bag.GetElements()
		if len(elements) != 1 || findCountInSlice(elements, "apple") != 3 {
			t.Errorf("Expected frequency 3 for 'apple', got %v", elements)
		}
	})
}

func TestBag_Remove(t *testing.T) {
	t.Run("decrement frequency", func(t *testing.T) {
		bag := InitBag[string](10)

		bag.Put("apple")
		bag.Put("apple")
		bag.Remove("apple")

		if bag.Size() != 1 {
			t.Errorf("Expected unique size 1 after removing one instance, got %d", bag.Size())
		}

		elements := bag.GetElements()
		if len(elements) != 1 || findCountInSlice(elements, "apple") != 1 {
			t.Errorf("Expected remaining frequency 1 for 'apple', got %v", elements)
		}
	})

	t.Run("complete removal and cluster rehash", func(t *testing.T) {
		bag := InitBag[string](5)

		bag.Put("A")
		bag.Put("B")
		bag.Put("B")
		bag.Put("C")

		bag.Remove("A")
		bag.Remove("C")

		if bag.Size() != 1 {
			t.Errorf("Expected unique size 1 after full removals, got %d", bag.Size())
		}

		elements := bag.GetElements()
		if len(elements) != 1 || findCountInSlice(elements, "B") != 2 {
			t.Errorf("Expected frequency 2 for 'B' to be preserved after rehash, got %v", elements)
		}
	})

	t.Run("remove non-existing element", func(t *testing.T) {
		bag := InitBag[string](10)

		bag.Put("apple")
		bag.Remove("banana")

		if bag.Size() != 1 {
			t.Errorf("Expected size to remain 1 when removing missing element, got %d", bag.Size())
		}
	})
}

func TestBag_GetElements(t *testing.T) {
	t.Run("empty bag", func(t *testing.T) {
		bag := InitBag[string](10)
		elements := bag.GetElements()

		if len(elements) != 0 {
			t.Errorf("Expected empty slice for empty bag, got length %d", len(elements))
		}
	})

	t.Run("populated bag with frequencies", func(t *testing.T) {
		bag := InitBag[string](10)

		bag.Put("red")
		bag.Put("blue")
		bag.Put("red")
		bag.Put("green")
		bag.Put("red")

		elements := bag.GetElements()
		if len(elements) != 3 {
			t.Errorf("Expected 3 unique elements in slice, got %d", len(elements))
		}
		if findCountInSlice(elements, "red") != 3 || findCountInSlice(elements, "blue") != 1 {
			t.Errorf("Expected frequencies red=3 and blue=1, got slice %v", elements)
		}
		if findCountInSlice(elements, "green") != 1 {
			t.Errorf("Expected frequency green=1, got slice %v", elements)
		}
	})
}


