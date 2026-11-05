package main

import (
	"testing"
)

func TestOrderedList_RemoveDuplicates(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		l.RemoveDuplicates()

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", l.Count())
		}
	})

	t.Run("Single element list", func(t *testing.T) {
		l := &OrderedList[int]{}
		l.Add(1)

		l.RemoveDuplicates()

		if l.Count() != 1 {
			t.Errorf("Expected 1 element, got %d", l.Count())
		}
		if l.head.value != 1 || l.tail.value != 1 {
			t.Errorf("Expected head and tail to be 1")
		}
	})

	t.Run("No duplicates ascending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		l.RemoveDuplicates()

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}

		validateListIntegrity(t, l)
	})

	t.Run("Duplicates ascending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(3)

		l.RemoveDuplicates()

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}
		if l.head.value != 1 || l.tail.value != 3 {
			t.Errorf("Expected head 1 and tail 3, got %d and %d", l.head.value, l.tail.value)
		}

		validateListIntegrity(t, l)
	})

	t.Run("Duplicates descending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}
		l.Add(3)
		l.Add(3)
		l.Add(2)
		l.Add(1)
		l.Add(1)

		l.RemoveDuplicates()

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}
		if l.head.value != 3 || l.tail.value != 1 {
			t.Errorf("Expected head 3 and tail 1, got %d and %d", l.head.value, l.tail.value)
		}

		validateListIntegrity(t, l)
	})

	t.Run("Multiple consecutive duplicates", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(1)
		l.Add(2)
		l.Add(2)

		l.RemoveDuplicates()

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", l.Count())
		}
		if l.head.value != 1 || l.tail.value != 2 {
			t.Errorf("Expected head 1 and tail 2")
		}

		validateListIntegrity(t, l)
	})

	t.Run("All elements are duplicates", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(5)
		l.Add(5)
		l.Add(5)

		l.RemoveDuplicates()

		if l.Count() != 1 {
			t.Errorf("Expected 1 element, got %d", l.Count())
		}
		if l.head.value != 5 || l.tail.value != 5 {
			t.Errorf("Expected head and tail to be 5")
		}

		validateListIntegrity(t, l)
	})

	t.Run("String duplicates", func(t *testing.T) {
		l := &OrderedList[string]{ascending: true}
		l.Add("a")
		l.Add("a")
		l.Add("b")

		l.RemoveDuplicates()

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", l.Count())
		}

		validateListIntegrity(t, l)
	})
}
