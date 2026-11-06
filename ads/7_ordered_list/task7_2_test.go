package main

import (
	"constraints"
	"math/rand"
	"sort"
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
		l := &OrderedList[int]{_ascending: true}
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
		l := &OrderedList[int]{_ascending: true}
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
		l := &OrderedList[int]{_ascending: false}
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
		l := &OrderedList[int]{_ascending: true}
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
		l := &OrderedList[int]{_ascending: true}
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
		l := &OrderedList[string]{_ascending: true}
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
		list1 := &OrderedList[int]{_ascending: true}
		list2 := &OrderedList[int]{_ascending: true}

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
		list1 := &OrderedList[int]{_ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{_ascending: true}

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
		list1 := &OrderedList[int]{_ascending: false}
		list1.Add(5)
		list1.Add(3)
		list1.Add(1)

		list2 := &OrderedList[int]{_ascending: false}

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
		list1 := &OrderedList[int]{_ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{_ascending: true}
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
		list1 := &OrderedList[int]{_ascending: false}
		list1.Add(6)
		list1.Add(4)
		list1.Add(2)

		list2 := &OrderedList[int]{_ascending: false}
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
		list1 := &OrderedList[int]{_ascending: true}
		list1.Add(1)
		list1.Add(3)
		list1.Add(5)

		list2 := &OrderedList[int]{_ascending: true}
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
		list1 := &OrderedList[int]{_ascending: true}
		list1.Add(1)
		list1.Add(2)

		list2 := &OrderedList[int]{_ascending: false}
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
		list1 := &OrderedList[int]{_ascending: true}
		list1.Add(5)

		list2 := &OrderedList[int]{_ascending: true}
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
		list1 := &OrderedList[string]{_ascending: true}
		list1.Add("apple")
		list1.Add("cherry")

		list2 := &OrderedList[string]{_ascending: true}
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
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)

		sub := &OrderedList[int]{_ascending: true}

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected empty list to be a sublist of any list")
		}
	})

	t.Run("Empty main list and empty sublist", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		sub := &OrderedList[int]{_ascending: true}

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected empty list to contain empty sublist")
		}
	})

	t.Run("Sublist at the beginning", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(4)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(1)
		sub.Add(2)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [1, 2] to be found at the beginning")
		}
	})

	t.Run("Sublist in the middle", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(10)
		l.Add(20)
		l.Add(30)
		l.Add(40)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(20)
		sub.Add(30)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [20, 30] to be found in the middle")
		}
	})

	t.Run("Sublist at the end", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(2)
		sub.Add(3)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 3] to be found at the end")
		}
	})

	t.Run("Sublist equals main list", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(5)
		l.Add(6)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(5)
		sub.Add(6)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected list to contain itself as a sublist")
		}
	})

	t.Run("Sublist not found", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(3)
		l.Add(5)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(2)
		sub.Add(3)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 3] not to be found in [1, 3, 5]")
		}
	})

	t.Run("Sublist larger than main list", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(1)
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false when sublist is larger than main list")
		}
	})

	t.Run("Mismatched ordering", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)

		sub := &OrderedList[int]{_ascending: false}
		sub.Add(1)
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false for lists with different ordering")
		}
	})

	t.Run("Single element sublist found", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(10)
		l.Add(20)
		l.Add(30)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(20)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [20] to be found")
		}
	})

	t.Run("String lists", func(t *testing.T) {
		l := &OrderedList[string]{_ascending: true}
		l.Add("apple")
		l.Add("banana")
		l.Add("cherry")

		sub := &OrderedList[string]{_ascending: true}
		sub.Add("banana")
		sub.Add("cherry")

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected ['banana', 'cherry'] to be found")
		}
	})

	t.Run("Sublist with duplicates in main list", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(3)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(2)
		sub.Add(2)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [2, 2] to be found in [1, 2, 2, 2, 3]")
		}
	})

	t.Run("Single element sublist not found", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(3)
		l.Add(5)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(2)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected [2] not to be found in [1, 3, 5]")
		}
	})

	t.Run("Sublist larger than main list (empty main)", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(1)

		if l.ContainsSublist(sub) {
			t.Errorf("Expected false when main list is empty and sublist is not")
		}
	})

	t.Run("Complex overlap with duplicates", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(3)
		l.Add(4)

		sub := &OrderedList[int]{_ascending: true}
		sub.Add(1)
		sub.Add(2)
		sub.Add(3)

		if !l.ContainsSublist(sub) {
			t.Errorf("Expected [1, 2, 3] to be found")
		}
	})
}

func TestOrderedList_MostFrequentVal(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		val, err := l.MostFrequentVal()

		if err != ErrEmptyList {
			t.Errorf("Expected ErrEmptyList, got %v", err)
		}

		var zero int
		if val != zero {
			t.Errorf("Expected zero value, got %d", val)
		}
	})

	t.Run("Single element", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(5)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected 5, got %d", val)
		}
	})

	t.Run("All elements unique ascending", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		// первый элемент должен вернуться, так как все встречаются одинаково
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})

	t.Run("All elements unique descending", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: false}
		l.Add(3)
		l.Add(2)
		l.Add(1)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 3 {
			t.Errorf("Expected 3, got %d", val)
		}
	})

	t.Run("One element appears most frequently ascending", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(3)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 2 {
			t.Errorf("Expected 2, got %d", val)
		}
	})

	t.Run("One element appears most frequently descending", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: false}
		l.Add(3)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(1)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 2 {
			t.Errorf("Expected 2, got %d", val)
		}
	})

	t.Run("Multiple elements with same max frequency", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(2)
		l.Add(2)
		l.Add(3)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		// должен вернуть первый элемент с максимальной частотой
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})

	t.Run("All elements are the same", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(5)
		l.Add(5)
		l.Add(5)
		l.Add(5)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 5 {
			t.Errorf("Expected 5, got %d", val)
		}
	})

	t.Run("Two elements with different frequencies", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(1)
		l.Add(1)
		l.Add(2)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 1 {
			t.Errorf("Expected 1, got %d", val)
		}
	})

	t.Run("Last element is most frequent", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(3)
		l.Add(3)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 3 {
			t.Errorf("Expected 3, got %d", val)
		}
	})

	t.Run("String list with duplicates", func(t *testing.T) {
		l := &OrderedList[string]{_ascending: true}
		l.Add("apple")
		l.Add("banana")
		l.Add("banana")
		l.Add("banana")
		l.Add("cherry")

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != "banana" {
			t.Errorf("Expected 'banana', got '%s'", val)
		}
	})

	t.Run("String list all unique", func(t *testing.T) {
		l := &OrderedList[string]{_ascending: true}
		l.Add("apple")
		l.Add("banana")
		l.Add("cherry")

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != "apple" {
			t.Errorf("Expected 'apple', got '%s'", val)
		}
	})

	t.Run("Large frequency difference", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(2)
		l.Add(3)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 2 {
			t.Errorf("Expected 2, got %d", val)
		}
	})

	t.Run("Two elements only both same", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(7)
		l.Add(7)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 7 {
			t.Errorf("Expected 7, got %d", val)
		}
	})

	t.Run("Two elements only different", func(t *testing.T) {
		l := &OrderedList[int]{_ascending: true}
		l.Add(7)
		l.Add(8)

		val, err := l.MostFrequentVal()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if val != 7 {
			t.Errorf("Expected 7, got %d", val)
		}
	})
}

func validateSkipListIntegrity[T constraints.Ordered](t *testing.T, l *SkipList[T]) {
	t.Helper()

	// проверка count и порядка на Level 0
	count := 0
	var prev *SkipNode[T]

	for node := l.head.next[0]; node != nil; node = node.next[0] {
		count++

		if prev != nil {
			if l._ascending && l.Compare(prev.value, node.value) > 0 {
				t.Errorf("Ordering violated (asc): %v > %v", prev.value, node.value)
			}
			if !l._ascending && l.Compare(prev.value, node.value) < 0 {
				t.Errorf("Ordering violated (desc): %v < %v", prev.value, node.value)
			}
		}

		prev = node
	}

	if count != l.count {
		t.Errorf("Count mismatch: expected %d, got %d", l.count, count)
	}

	// проверка span на всех уровнях
	for lev := range l.level {
		curr := l.head
		spanSum := 0

		for {
			if curr.span[lev] <= 0 {
				t.Errorf("Invalid span <= 0 at level %d", lev)
			}

			spanSum += curr.span[lev]

			if curr.next[lev] == nil {
				break
			}

			curr = curr.next[lev]
		}

		if spanSum != count+1 {
			t.Errorf("Span sum mismatch at level %d: expected %d, got %d", lev, count+1, spanSum)
		}
	}
}

type referenceList[T constraints.Ordered] struct {
	items      []T
	_ascending bool
}

func (r *referenceList[T]) add(less func(T, T) int, val T) {
	r.items = append(r.items, val)

	sort.Slice(r.items, func(i, j int) bool {
		compare := less(r.items[i], r.items[j])

		if r._ascending {
			return compare < 0
		}

		return compare > 0
	})
}

func (r *referenceList[T]) delete(less func(T, T) int, val T) bool {
	for i, item := range r.items {
		if less(item, val) == 0 {
			r.items = append(r.items[:i], r.items[i+1:]...)

			return true
		}
	}

	return false
}

func validateSkipListWithValues[T constraints.Ordered](t *testing.T, l *SkipList[T], expectedItems []T) {
	t.Helper()

	var actualItems []T
	for node := l.head.next[0]; node != nil; node = node.next[0] {
		actualItems = append(actualItems, node.value)
	}

	if len(actualItems) != l.count {
		t.Errorf("Level 0 count expecetd: %d, got %d", l.count, len(actualItems))
	}

	if len(expectedItems) != len(actualItems) {
		t.Fatalf("Expected length: %d, got: %d", len(expectedItems), len(actualItems))
	}

	for i := range expectedItems {
		if l.Compare(expectedItems[i], actualItems[i]) != 0 {
			t.Errorf("Mismatch at index %d: expected %v, got %v", i, expectedItems[i], actualItems[i])
		}
	}

	// 2. Строгая проверка корректности span на абсолютно всех уровнях
	for lev := range l.level {
		curr := l.head
		spanSum := 0
		nodeCount := 0

		for {
			if curr.span[lev] <= 0 {
				t.Errorf("Incorrect or negative span (%d) at level %d", curr.span[lev], lev)
			}

			spanSum += curr.span[lev]

			if curr.next[lev] == nil {
				break
			}

			// сумма span шагов должна совпадать с физическим расстоянием на базовом уровне 0
			// Идем по Level 0 и считаем шаги до следующего узла текущего уровня
			stepsOnLevel0 := 0
			for l0 := curr; l0 != curr.next[lev]; l0 = l0.next[0] {
				stepsOnLevel0++
			}
			if curr.span[lev] != stepsOnLevel0 {
				t.Errorf("Wrong span calculation at level %d of a node with value %v: expected %d, got at Level 0 distance %d",
					lev, curr.value, curr.span[lev], stepsOnLevel0)
			}

			curr = curr.next[lev]
			nodeCount++
		}

		// сумма всех span на любом уровне должна быть строго равна (количество элементов + 1)
		if spanSum != len(actualItems)+1 {
			t.Errorf("Wrong span sum at Level %d: expected %d, got %d", lev, len(actualItems)+1, spanSum)
		}
	}

	// 3. Сквозная проверка FindByIndex для каждого элемента
	for i, expVal := range expectedItems {
		gotVal, err := l.FindByIndex(i)

		if err != nil {
			t.Errorf("FindByIndex(%d) returns error: %v", i, err)
		}
		if l.Compare(expVal, gotVal) != 0 {
			t.Errorf("FindByIndex(%d) returns incorret value: expected %v, got %v", i, expVal, gotVal)
		}
	}
}
func TestSkipList_Add(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := NewSkipList[int](true)

		l.Add(5)

		if l.count != 1 {
			t.Errorf("Expected count 1, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})

	t.Run("Multiple elements ascending", func(t *testing.T) {
		l := NewSkipList[int](true)

		for _, v := range []int{5, 2, 8, 1, 9} {
			l.Add(v)
		}

		if l.count != 5 {
			t.Errorf("Expected count 5, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})

	t.Run("Multiple elements descending", func(t *testing.T) {
		l := NewSkipList[int](false)

		for _, v := range []int{5, 2, 8, 1, 9} {
			l.Add(v)
		}

		if l.count != 5 {
			t.Errorf("Expected count 5, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})

	t.Run("Duplicates", func(t *testing.T) {
		l := NewSkipList[int](true)

		l.Add(3)
		l.Add(3)
		l.Add(3)

		if l.count != 3 {
			t.Errorf("Expected count 3, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})
}

func TestSkipList_FindByIndex(t *testing.T) {
	t.Run("Valid indices", func(t *testing.T) {
		l := NewSkipList[int](true)
		for _, v := range []int{10, 20, 30, 40, 50} {
			l.Add(v)
		}

		tests := []struct {
			index    int
			expected int
		}{
			{0, 10},
			{2, 30},
			{4, 50},
		}

		for _, tc := range tests {
			val, err := l.FindByIndex(tc.index)

			if err != nil {
				t.Errorf("Index %d: expected no error, got %v", tc.index, err)
			}
			if val != tc.expected {
				t.Errorf("Index %d: expected %d, got %d", tc.index, tc.expected, val)
			}
		}
	})

	t.Run("Out of bounds", func(t *testing.T) {
		l := NewSkipList[int](true)
		l.Add(10)

		_, err := l.FindByIndex(-1)
		if err == nil {
			t.Errorf("Expected error for index -1, got nil")
		}

		_, err = l.FindByIndex(5)
		if err == nil {
			t.Errorf("Expected error for index 5, got nil")
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		l := NewSkipList[int](true)

		// empty list
		_, err := l.FindByIndex(0)
		if err != ErrInvalidIndex {
			t.Errorf("Expected ErrInvalidIndex, got: %v", err)
		}

		l.Add(10)
		l.Add(20)

		// invalid incides
		_, err = l.FindByIndex(2)
		if err != ErrInvalidIndex {
			t.Errorf("Expected ErrInvalinIndex, got: %v", err)
		}

		_, err = l.FindByIndex(-1)
		if err != ErrInvalidIndex {
			t.Errorf("Expected ErrInvalinIndex, got: %v", err)
		}
	})
}

func TestSkipList_Delete(t *testing.T) {
	t.Run("Delete existing element", func(t *testing.T) {
		l := NewSkipList[int](true)
		l.Add(10)
		l.Add(20)
		l.Add(30)

		deleted := l.Delete(20)

		if !deleted {
			t.Errorf("Expected true, got false")
		}
		if l.count != 2 {
			t.Errorf("Expected count 2, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})

	t.Run("Delete non-existing element", func(t *testing.T) {
		l := NewSkipList[int](true)
		l.Add(10)
		l.Add(30)

		deleted := l.Delete(20)

		if deleted {
			t.Errorf("Expected false, got true")
		}
		if l.count != 2 {
			t.Errorf("Expected count 2, got %d", l.count)
		}
	})

	t.Run("Delete makes list empty", func(t *testing.T) {
		l := NewSkipList[int](true)
		l.Add(5)

		l.Delete(5)

		if l.count != 0 {
			t.Errorf("Expected count 0, got %d", l.count)
		}
		if l.level != 1 {
			t.Errorf("Expected level 1, got %d", l.level)
		}
	})

	t.Run("Delete duplicates", func(t *testing.T) {
		l := NewSkipList[int](true)
		l.Add(5)
		l.Add(5)
		l.Add(5)

		l.Delete(5)

		if l.count != 2 {
			t.Errorf("Expected count 2, got %d", l.count)
		}

		validateSkipListIntegrity(t, l)
	})
}

func TestSkipList_StressRandomOperations(t *testing.T) {
	t.Run("Ascending", func(t *testing.T) {
		seed := int64(42) // фиксируем seed для воспроизводимости в случае падения
		rng := rand.New(rand.NewSource(seed))

		l := NewSkipList[int](true)
		ref := &referenceList[int]{_ascending: true}

		// 500 случайных операций (Add / Delete / validateSkipListWithValues)
		for range 500 {
			op := rng.Intn(10)

			if op < 6 { // 60% вероятность добавления
				val := rng.Intn(200) // могут быть дубликаты

				l.Add(val)
				ref.add(l.Compare, val)

			} else if len(ref.items) > 0 { // 40% вероятность удаления
				// удаляется либо случайный существующий, либо вообще случайное число
				var val int

				if rng.Float64() < 0.8 {
					val = ref.items[rng.Intn(len(ref.items))]
				} else {
					val = rng.Intn(200)
				}

				lDeleted := l.Delete(val)
				refDeleted := ref.delete(l.Compare, val)

				if lDeleted != refDeleted {
					t.Fatalf("Mismatch on Delete(%d): SkipList=%v, Slice=%v", val, lDeleted, refDeleted)
				}
			}

			// проверка структуры на каждой итерации
			validateSkipListWithValues(t, l, ref.items)
		}
	})

	t.Run("Descending", func(t *testing.T) {
		rng := rand.New(rand.NewSource(1337))

		l := NewSkipList[int](false)
		ref := &referenceList[int]{_ascending: false}

		for range 500 {
			op := rng.Intn(10)

			if op < 6 {
				val := rng.Intn(200)

				l.Add(val)
				ref.add(l.Compare, val)
			} else if len(ref.items) > 0 {
				val := ref.items[rng.Intn(len(ref.items))]

				lDeleted := l.Delete(val)
				refDeleted := ref.delete(l.Compare, val)

				if lDeleted != refDeleted {
					t.Fatalf("Mismatch on Delete(%d): SkipList=%v, Slice=%v", val, lDeleted, refDeleted)
				}
			}

			validateSkipListWithValues(t, l, ref.items)
		}
	})
}


