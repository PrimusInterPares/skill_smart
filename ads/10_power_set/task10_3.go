package main

import (
	"strconv"
	"testing"
	"time"
)

func TestPowerSet_Put(t *testing.T) {
	t.Run("basic insertion", func(t *testing.T) {
		set := Init[string]()

		set.Put("apple")

		if set.Size() != 1 {
			t.Errorf("Expected size 1 after putting 'apple', got %d", set.Size())
		}
		if !set.Get("apple") {
			t.Errorf("Expected 'apple' to be present in the set")
		}
	})

	t.Run("duplicate insertion", func(t *testing.T) {
		set := Init[string]()

		set.Put("apple")
		set.Put("apple")

		if set.Size() != 1 {
			t.Errorf("Expected size to remain 1 when putting a duplicate element, got %d", set.Size())
		}
	})
}

func TestPowerSet_Get(t *testing.T) {
	t.Run("non-existing element", func(t *testing.T) {
		set := Init[string]()

		if set.Get("banana") {
			t.Errorf("Expected Get('banana') to be false for an empty set")
		}
	})

	t.Run("existing element", func(t *testing.T) {
		set := Init[string]()

		set.Put("banana")

		if !set.Get("banana") {
			t.Errorf("Expected Get('banana') to be true after insertion")
		}
	})
}

func TestPowerSet_Remove(t *testing.T) {
	t.Run("existing element", func(t *testing.T) {
		set := Init[string]()

		set.Put("cherry")

		if !set.Remove("cherry") {
			t.Errorf("Expected Remove('cherry') to return true for an existing element")
		}
		if set.Get("cherry") {
			t.Errorf("Expected 'cherry' to be removed from the set")
		}
		if set.Size() != 0 {
			t.Errorf("Expected size 0 after removal, got %d", set.Size())
		}
	})

	t.Run("missing element", func(t *testing.T) {
		set := Init[string]()

		if set.Remove("missing") {
			t.Errorf("Expected Remove('missing') to return false for a non-existing element")
		}
	})
}

func TestPowerSet_Intersection(t *testing.T) {
	t.Run("non-empty intersection", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("A")
		set1.Put("B")
		set2.Put("B")
		set2.Put("C")

		inter := set1.Intersection(set2)
		if inter.Size() != 1 || !inter.Get("B") {
			t.Errorf("Expected intersection to contain only 'B', got size %d", inter.Size())
		}
	})

	t.Run("empty intersection", func(t *testing.T) {
		set1 := Init[string]()

		set1.Put("A")
		set1.Put("B")

		set3 := Init[string]()

		set3.Put("Z")

		emptyInter := set1.Intersection(set3)
		if emptyInter.Size() != 0 {
			t.Errorf("Expected intersection to be empty, got size %d", emptyInter.Size())
		}
	})
}

func TestPowerSet_Union(t *testing.T) {
	t.Run("both non-empty", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("A")
		set2.Put("B")

		union1 := set1.Union(set2)
		if union1.Size() != 2 || !union1.Get("A") || !union1.Get("B") {
			t.Errorf("Expected union to contain 'A' and 'B', got size %d", union1.Size())
		}
	})

	t.Run("one empty", func(t *testing.T) {
		set1 := Init[string]()

		set1.Put("A")

		emptySet := Init[string]()

		union2 := set1.Union(emptySet)
		if union2.Size() != 1 || !union2.Get("A") {
			t.Errorf("Expected union with empty set to return an identical copy, got size %d", union2.Size())
		}
	})
}

func TestPowerSet_Difference(t *testing.T) {
	t.Run("non-empty difference", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("A")
		set1.Put("B")
		set2.Put("B")

		diff := set1.Difference(set2)
		if diff.Size() != 1 || !diff.Get("A") || diff.Get("B") {
			t.Errorf("Expected difference to contain only 'A', got size %d", diff.Size())
		}
	})

	t.Run("empty difference", func(t *testing.T) {
		set1 := Init[string]()

		set1.Put("A")
		set1.Put("B")
		set3 := Init[string]()

		set3.Put("A")
		set3.Put("B")

		emptyDiff := set1.Difference(set3)
		if emptyDiff.Size() != 0 {
			t.Errorf("Expected difference to be empty, got size %d", emptyDiff.Size())
		}
	})
}

func TestPowerSet_IsSubset(t *testing.T) {
	t.Run("valid subset", func(t *testing.T) {
		setMain := Init[string]()

		setMain.Put("1")
		setMain.Put("2")
		setMain.Put("3")

		sub1 := Init[string]()

		sub1.Put("1")
		sub1.Put("2")

		if !setMain.IsSubset(sub1) {
			t.Errorf("Expected sub1 to be a subset of setMain")
		}
	})

	t.Run("not subset (larger)", func(t *testing.T) {
		setMain := Init[string]()

		setMain.Put("1")
		setMain.Put("2")
		setMain.Put("3")

		sub2 := Init[string]()

		sub2.Put("1")
		sub2.Put("2")
		sub2.Put("3")
		sub2.Put("4")

		if setMain.IsSubset(sub2) {
			t.Errorf("Expected setMain to NOT be a subset of a larger sub2")
		}
	})

	t.Run("not subset (partial overlap)", func(t *testing.T) {
		setMain := Init[string]()

		setMain.Put("1")
		setMain.Put("2")
		setMain.Put("3")

		sub3 := Init[string]()

		sub3.Put("1")
		sub3.Put("5")
		if setMain.IsSubset(sub3) {
			t.Errorf("Expected sub3 to NOT be a subset due to missing element '5'")
		}
	})
}

func TestPowerSet_Equals(t *testing.T) {
	t.Run("equal sets", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("X")
		set1.Put("Y")
		set2.Put("Y")
		set2.Put("X")

		if !set1.Equals(set2) {
			t.Errorf("Expected sets to be equal regardless of insertion order")
		}
	})

	t.Run("unequal sets", func(t *testing.T) {
		set1 := Init[string]()
		set2 := Init[string]()

		set1.Put("X")
		set1.Put("Y")
		set2.Put("Y")
		set2.Put("X")
		set2.Put("Z")

		if set1.Equals(set2) {
			t.Errorf("Expected sets with different item counts to NOT be equal")
		}
	})
}

func TestPowerSet_Integration(t *testing.T) {
	t.Run("cross-functional pipeline", func(t *testing.T) {
		setA := Init[string]()
		setB := Init[string]()

		setA.Put("red")
		setA.Put("green")
		setA.Put("blue")
		setB.Put("green")
		setB.Put("blue")
		setB.Put("yellow")

		inter := setA.Intersection(setB)
		diff := setA.Difference(inter)
		if diff.Size() != 1 || !diff.Get("red") {
			t.Errorf("Integration failed at Difference step, expected only 'red', got size %d", diff.Size())
		}

		diff.Remove("red")
		if !setB.IsSubset(diff) {
			t.Errorf("An empty set must always be a subset of any other set")
		}
	})
}

func TestPowerSet_Performance(t *testing.T) {
	set1 := Init[string]()
	set2 := Init[string]()
	elementsCount := 15000

	t.Run("heavy insertions", func(t *testing.T) {
		for i := range elementsCount {
			set1.Put("item_" + strconv.Itoa(i))
			set2.Put("item_" + strconv.Itoa(i+5000))
		}
	})

	t.Run("sequential mass operations", func(t *testing.T) {
		startOps := time.Now()
		_ = set1.Intersection(set2)
		_ = set1.Union(set2)
		_ = set1.Difference(set2)
		_ = set1.IsSubset(set2)
		_ = set1.Equals(set2)
		durationOps := time.Since(startOps)

		t.Logf("Bulk algebraic operations took: %v", durationOps)

		if durationOps > 1*time.Second {
			t.Errorf("Algebraic operations took too long: %v (expected < 1s)", durationOps)
		}
	})

	t.Run("mass removals", func(t *testing.T) {
		startRemove := time.Now()

		for i := 0; i < elementsCount; i += 3 {
			set1.Remove("item_" + strconv.Itoa(i))
		}

		durationRemove := time.Since(startRemove)

		t.Logf("Mass removal of 5000 elements took: %v", durationRemove)
		if durationRemove > 1*time.Second {
			t.Errorf("Mass removal operation took too long: %v (expected < 1s)", durationRemove)
		}
	})

	t.Run("data integrity", func(t *testing.T) {
		for i := range elementsCount {
			exists := set1.Get("item_" + strconv.Itoa(i))

			if i%3 == 0 && exists {
				t.Errorf("Element item_%d should be deleted, but it was found", i)
			}
			if i%3 != 0 && !exists {
				t.Errorf("Element item_%d was mistakenly lost due to a broken collision chain!", i)
			}
		}
	})
}


