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

func TestMergeOrderedLists(t *testing.T) {
	t.Run("Empty lists", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list2 := &OrderedList[int]{ascending: true}

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", merged.Count())
		}
		validateListIntegrity(t, merged)
	})

	t.Run("One empty list ascending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{ascending: true}

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", merged.Count())
		}
		if merged.head.value != 1 || merged.tail.value != 5 {
			t.Errorf("Expected head=1, tail=5")
		}
		validateListIntegrity(t, merged)
	})

	t.Run("One empty list descending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: false}
		list1.Add(5)
		list1.Add(3)
		list1.Add(1)

		list2 := &OrderedList[int]{ascending: false}

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", merged.Count())
		}
		if merged.head.value != 5 || merged.tail.value != 1 {
			t.Errorf("Expected head=5, tail=1")
		}
		validateListIntegrity(t, merged)
	})

	t.Run("Both lists with elements ascending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{ascending: true}
		list2.Add(2)
		list2.Add(4)
		list2.Add(6)

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 6 {
			t.Errorf("Expected 6 elements, got %d", merged.Count())
		}

		// Check order: 1, 2, 3, 4, 5, 6
		expected := []int{1, 2, 3, 4, 5, 6}
		node := merged.head
		for i, exp := range expected {
			if node == nil {
				t.Errorf("Expected element at position %d", i)
				break
			}
			if node.value != exp {
				t.Errorf("Position %d: expected %d, got %d", i, exp, node.value)
			}
			node = node.next
		}
		validateListIntegrity(t, merged)
	})

	t.Run("Both lists with elements descending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: false}
		list1.Add(6)
		list1.Add(4)
		list1.Add(2)

		list2 := &OrderedList[int]{ascending: false}
		list2.Add(5)
		list2.Add(3)
		list2.Add(1)

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 6 {
			t.Errorf("Expected 6 elements, got %d", merged.Count())
		}

		// Check order: 6, 5, 4, 3, 2, 1
		expected := []int{6, 5, 4, 3, 2, 1}
		node := merged.head
		for i, exp := range expected {
			if node == nil {
				t.Errorf("Expected element at position %d", i)
				break
			}
			if node.value != exp {
				t.Errorf("Position %d: expected %d, got %d", i, exp, node.value)
			}
			node = node.next
		}
		validateListIntegrity(t, merged)
	})

	t.Run("Overlapping values ascending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{ascending: true}
		list2.Add(2)
		list2.Add(3)
		list2.Add(4)

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 6 {
			t.Errorf("Expected 6 elements, got %d", merged.Count())
		}

		// Check order: 1, 2, 3, 3, 4, 5
		expected := []int{1, 2, 3, 3, 4, 5}
		node := merged.head
		for i, exp := range expected {
			if node == nil {
				t.Errorf("Expected element at position %d", i)
				break
			}
			if node.value != exp {
				t.Errorf("Position %d: expected %d, got %d", i, exp, node.value)
			}
			node = node.next
		}
		validateListIntegrity(t, merged)
	})

	t.Run("Mismatched ordering", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list1.Add(1)
		list1.Add(2)

		list2 := &OrderedList[int]{ascending: false}
		list2.Add(3)
		list2.Add(4)

		merged, err := MergeOrderedLists(list1, list2)
		if err != ErrOrderingMismatch {
			t.Errorf("Expected ErrOrderingMismatch, got %v", err)
		}
		if merged.Count() != 0 {
			t.Errorf("Expected 0 elements for error case, got %d", merged.Count())
		}
	})

	t.Run("Single element lists ascending", func(t *testing.T) {
		list1 := &OrderedList[int]{ascending: true}
		list1.Add(5)

		list2 := &OrderedList[int]{ascending: true}
		list2.Add(3)

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", merged.Count())
		}
		if merged.head.value != 3 || merged.tail.value != 5 {
			t.Errorf("Expected head=3, tail=5")
		}
		validateListIntegrity(t, merged)
	})

	t.Run("String lists ascending", func(t *testing.T) {
		list1 := &OrderedList[string]{ascending: true}
		list1.Add("apple")
		list1.Add("cherry")

		list2 := &OrderedList[string]{ascending: true}
		list2.Add("banana")
		list2.Add("date")

		merged, err := MergeOrderedLists(list1, list2)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if merged.Count() != 4 {
			t.Errorf("Expected 4 elements, got %d", merged.Count())
		}

		expected := []string{"apple", "banana", "cherry", "date"}
		node := merged.head
		for i, exp := range expected {
			if node == nil {
				t.Errorf("Expected element at position %d", i)
				break
			}
			if node.value != exp {
				t.Errorf("Position %d: expected %s, got %s", i, exp, node.value)
			}
			node = node.next
		}
		validateListIntegrity(t, merged)
	})
}
