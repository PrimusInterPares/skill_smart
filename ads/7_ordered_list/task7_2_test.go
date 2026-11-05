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

func TestOrderedList_ContainsSublist(t *testing.T) {
	t.Run("Empty sublist in non-empty list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)

		sub := &OrderedList[int]{ascending: true}

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected empty list to be a sublist of any list")
		}
	})

	t.Run("Empty main list and empty sublist", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		sub := &OrderedList[int]{ascending: true}

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected empty list to contain empty sublist")
		}
	})

	t.Run("Sublist at the beginning", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(4)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(1)
		sub.Add(2)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [1, 2] to be found at the beginning")
		}
	})

	t.Run("Sublist in the middle", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(10)
		l.Add(20)
		l.Add(30)
		l.Add(40)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(20)
		sub.Add(30)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [20, 30] to be found in the middle")
		}
	})

	t.Run("Sublist at the end", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(2)
		sub.Add(3)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 3] to be found at the end")
		}
	})

	t.Run("Sublist equals main list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(5)
		l.Add(6)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(5)
		sub.Add(6)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected list to contain itself as a sublist")
		}
	})

	t.Run("Sublist not found", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(3)
		l.Add(5)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(2)
		sub.Add(3)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 3] not to be found in [1, 3, 5]")
		}
	})

	t.Run("Sublist larger than main list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(1)
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false when sublist is larger than main list")
		}
	})

	t.Run("Mismatched ordering", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)

		sub := &OrderedList[int]{ascending: false}
		sub.Add(1)
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false for lists with different ordering")
		}
	})

	t.Run("Single element sublist found", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(10)
		l.Add(20)
		l.Add(30)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(20)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [20] to be found")
		}
	})

	t.Run("String lists", func(t *testing.T) {
		l := &OrderedList[string]{ascending: true}
		l.Add("apple")
		l.Add("banana")
		l.Add("cherry")

		sub := &OrderedList[string]{ascending: true}
		sub.Add("banana")
		sub.Add("cherry")

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected ['banana', 'cherry'] to be found")
		}
	})

	t.Run("Sublist with duplicates in main list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(3)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(2)
		sub.Add(2)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 2] to be found in [1, 2, 2, 2, 3]")
		}
	})

	t.Run("Single element sublist not found", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(3)
		l.Add(5)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected [2] not to be found in [1, 3, 5]")
		}
	})

	t.Run("Sublist larger than main list (empty main)", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}

		sub := &OrderedList[int]{ascending: true}
		sub.Add(1)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false when main list is empty and sublist is not")
		}
	})

	t.Run("Complex overlap with duplicates", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(3)
		l.Add(4)

		sub := &OrderedList[int]{ascending: true}
		sub.Add(1)
		sub.Add(2)
		sub.Add(3)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [1, 2, 3] to be found")
		}
	})
}
