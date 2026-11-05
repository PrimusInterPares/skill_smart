package main

import (
	"constraints"
	"testing"
)

func validateListIntegrity[T constraints.Ordered](t *testing.T, l *OrderedList[T]) {
	t.Helper()

	if l.Count() == 0 {
		if l.head != nil || l.tail != nil {
			t.Errorf("Empty list must have nil head and tail")
		}
		return
	}

	if l.head.prev != nil {
		t.Errorf("Head's prev must be nil")
	}
	if l.tail.next != nil {
		t.Errorf("Tail's next must be nil")
	}

	node := l.head
	count := 0
	for node != nil {
		count++
		if node.next != nil && node.next.prev != node {
			t.Errorf("Broken link: node.next.prev != node (value: %v)", node.value)
		}
		if node.prev != nil && node.prev.next != node {
			t.Errorf("Broken link: node.prev.next != node (value: %v)", node.value)
		}
		node = node.next
	}

	if count != l.Count() {
		t.Errorf("Traversal count %d does not match l.Count() %d", count, l.Count())
	}
}

func TestOrderedList_Count(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", l.Count())
		}
	})

	t.Run("After additions", func(t *testing.T) {
		l := &OrderedList[int]{}
		l.Add(1)
		l.Add(2)

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", l.Count())
		}
	})
}

func TestOrderedList_Add(t *testing.T) {
	t.Run("Add to empty list", func(t *testing.T) {
		l := &OrderedList[int]{}
		l.Add(1)

		if l.Count() != 1 {
			t.Errorf("Expected 1 element, got %d", l.Count())
		}

		if l.head.value != 1 || l.tail.value != 1 {
			t.Errorf("Expected head and tail to be 1")
		}
	})

	t.Run("Add multiple elements into ascending list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}

		l.Add(3)
		l.Add(1)
		l.Add(2)

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}

		if l.head.value != 1 {
			t.Errorf("Expected head to be 1, got %d", l.head.value)
		}

		if l.tail.value != 3 {
			t.Errorf("Expected tail to be 3, got %d", l.tail.value)
		}

		validateListIntegrity(t, l)
	})

	t.Run("Add multiple elements into descending list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}

		l.Add(1)
		l.Add(3)
		l.Add(2)

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}

		if l.head.value != 3 {
			t.Errorf("Expected head to be 3, got %d", l.head.value)
		}

		if l.tail.value != 1 {
			t.Errorf("Expected tail to be 1, got %d", l.tail.value)
		}
	})

	t.Run("Add duplicate values ascending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}

		l.Add(2)
		l.Add(2)
		l.Add(2)

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}
		if l.head.value != 2 || l.tail.value != 2 {
			t.Errorf("Expected head and tail to be 2")
		}
	})

	t.Run("Add duplicate values descending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}

		l.Add(5)
		l.Add(5)
		l.Add(5)

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}
		if l.head.value != 5 || l.tail.value != 5 {
			t.Errorf("Expected head and tail to be 5")
		}
	})

	t.Run("Add smaller element to single-element ascending list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(5)

		l.Add(2)

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", l.Count())
		}
		if l.head.value != 2 {
			t.Errorf("Expected head to be 2, got %d", l.head.value)
		}
		if l.tail.value != 5 {
			t.Errorf("Expected tail to be 5, got %d", l.tail.value)
		}
	})
	t.Run("Add larger element to single-element ascending list", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(2)

		l.Add(5)

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements, got %d", l.Count())
		}
		if l.head.value != 2 {
			t.Errorf("Expected head to be 2, got %d", l.head.value)
		}
		if l.tail.value != 5 {
			t.Errorf("Expected tail to be 5, got %d", l.tail.value)
		}
	})
}

func TestOrderedList_Find(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		_, err := l.Find(1)

		if err != ErrEmptyList {
			t.Errorf("Expected ErrEmptyList, got %v", err)
		}
	})

	t.Run("Ascending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}

		l.Add(1)
		l.Add(3)
		l.Add(5)

		if node, err := l.Find(3); err != nil || node.value != 3 {
			t.Errorf("Expected 3, got %d, err %v", node.value, err)
		}
		if _, err := l.Find(4); err != ErrNodeNotFound {
			t.Errorf("Expected ErrNodeNotFound for 4, got %v", err)
		}
	})

	t.Run("Descending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}

		l.Add(5)
		l.Add(1)
		l.Add(3)

		if node, err := l.Find(3); err != nil || node.value != 3 {
			t.Errorf("Expected 3, got %d, err %v", node.value, err)
		}
		if _, err := l.Find(4); err != ErrNodeNotFound {
			t.Errorf("Expected ErrNodeNotFound for 4, got %v", err)
		}
	})

	t.Run("Find among duplicates ascending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(2)
		l.Add(2)
		l.Add(2)

		node, err := l.Find(2)
		if err != nil || node.value != 2 {
			t.Errorf("Expected to find 2, got %d, err %v", node.value, err)
		}
	})

	t.Run("Find among duplicates descending", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}
		l.Add(7)
		l.Add(7)
		l.Add(7)

		node, err := l.Find(7)
		if err != nil || node.value != 7 {
			t.Errorf("Expected to find 7, got %d, err %v", node.value, err)
		}
	})
}

func TestOrderedList_Delete(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		l.Delete(555)

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", l.Count())
		}
	})

	t.Run("List with one element", func(t *testing.T) {
		l := &OrderedList[string]{}
		l.Add("test")

		if l.Count() != 1 {
			t.Errorf("Expected 1 element in the list, got %d", l.Count())
		}

		l.Delete("test")

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements after delete, got %d", l.Count())
		}

		if l.head != nil || l.tail != nil {
			t.Errorf("Head and tail are expected to be nil and nil, got %v and %v", l.head, l.tail)
		}
	})

	t.Run("Delete head and tail", func(t *testing.T) {
		l := &OrderedList[int]{} // descending
		l.Add(1)
		l.Add(2)
		l.Add(3) // 3 (head) <-> 2 <-> 1 (tail)

		l.Delete(3)

		if l.Count() != 2 {
			t.Errorf("Expected 2 elements after head deletion, got %d", l.Count())
		}

		if l.head.value != 2 {
			t.Errorf("Expected new head value is 2, got %d", l.head.value)
		}

		l.Delete(1)

		if l.Count() != 1 {
			t.Errorf("Expected 1 elements after tail deletion, got %d", l.Count())
		}

		if l.tail.value != 2 {
			t.Errorf("Expected new tail value is 2, got %d", l.tail.value)
		}
	})

	t.Run("Delete middle element", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)
		l.Add(4)
		l.Add(5)

		l.Delete(3)

		if l.Count() != 4 {
			t.Errorf("Expected 4 elements, got %d", l.Count())
		}

		node, err := l.Find(2)
		if err != nil || node.next.value != 4 {
			t.Errorf("Expected 2 to be followed by 4")
		}

		validateListIntegrity(t, l)
	})

	t.Run("Delete non-existent element", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		l.Delete(99)

		if l.Count() != 3 {
			t.Errorf("Expected 3 elements, got %d", l.Count())
		}
	})

	t.Run("Delete all elements one by one", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}
		l.Add(1)
		l.Add(2)
		l.Add(3)

		l.Delete(2)
		l.Delete(1)
		l.Delete(3)

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", l.Count())
		}
		if l.head != nil || l.tail != nil {
			t.Errorf("Expected head and tail to be nil")
		}
	})
}

func TestOrderedList_Clear(t *testing.T) {
	t.Run("Clear list with elements and set ascending", func(t *testing.T) {
		l := &OrderedList[int]{}
		l.Add(1)
		l.Add(2)

		l.Clear(true)

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements after clear, got %d", l.Count())
		}

		if l.head != nil || l.tail != nil {
			t.Errorf("Expected head and tail to be nil after clear")
		}

		if !l.ascending {
			t.Errorf("Expected ascending to be true, got %v", l.ascending)
		}
	})

	t.Run("Clear list with elements and set descending", func(t *testing.T) {
		l := &OrderedList[int]{}
		l.Add(1)
		l.Add(2)

		l.Clear(false)

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements after clear, got %d", l.Count())
		}

		if l.head != nil || l.tail != nil {
			t.Errorf("Expected head and tail to be nil after clear")
		}

		if l.ascending {
			t.Errorf("Expected ascending to be false, got %v", l.ascending)
		}
	})

	t.Run("Clear already empty list", func(t *testing.T) {
		l := &OrderedList[int]{}

		l.Clear(true)

		if l.Count() != 0 {
			t.Errorf("Expected 0 elements, got %d", l.Count())
		}
		if l.head != nil || l.tail != nil {
			t.Errorf("Expected head and tail to be nil")
		}
		if !l.ascending {
			t.Errorf("Expected ascending to be true")
		}
	})
}

func TestOrderedList_Compare(t *testing.T) {
	t.Run("Ascending is true", func(t *testing.T) {
		l := &OrderedList[int]{ascending: true}

		if l.Compare(1, 2) != -1 {
			t.Errorf("Expected -1 for 1 < 2, got %d", l.Compare(1, 2))
		}

		if l.Compare(2, 1) != 1 {
			t.Errorf("Expected 1 for 2 > 1, got %d", l.Compare(2, 1))
		}

		if l.Compare(1, 1) != 0 {
			t.Errorf("Expected 0 for 1 == 1, got %d", l.Compare(1, 1))
		}
	})

	t.Run("Ascending is false", func(t *testing.T) {
		l := &OrderedList[int]{ascending: false}

		if l.Compare(1, 2) != 1 {
			t.Errorf("Expected 1 for 1 < 2 in descending, got %d", l.Compare(1, 2))
		}

		if l.Compare(2, 1) != -1 {
			t.Errorf("Expected -1 for 2 > 1 in descending, got %d", l.Compare(2, 1))
		}

		if l.Compare(1, 1) != 0 {
			t.Errorf("Expected 0 for 1 == 1 in descending, got %d", l.Compare(1, 1))
		}
	})
}
