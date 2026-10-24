package main

import(
	"testing"
)

func TestReverseList(t *testing.T) {
	t.Run("One-element list", func(t *testing.T) {
		l := createList([]int{1})
		reversed := ReverseList(l)

		if reversed.head != reversed.tail {
			t.Fatalf("Unexpected changes: head and tail are different")
		}
	})

	t.Run("Multiple elements list", func(t *testing.T) {
		l := createList([]int{1, 2, 3})
		reversed := ReverseList(l)

		rFirst := reversed.head
		rSecond := rFirst.next
		rThird := rSecond.next
		if rFirst.value != 3 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 3)
		}

		if rFirst.next != rSecond {
			t.Errorf("rFirst next is not rSecond")
		}

		if rFirst.prev != nil {
			t.Errorf("rFirst prev is not nil")
		}
		
		if rSecond.value != 2 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 2)
		}

		if rSecond.next != rThird {
			t.Errorf("rSecond next is not rThird")
		}

		if rSecond.prev != rFirst {
			t.Errorf("rSecond prev is not rFirst")
		}

		if rThird.value != 1 {
			t.Errorf("rThird value is incorrect: %d, expected: %d", rThird.value, 1)
		}

		if rThird.next != nil {
			t.Errorf("rThird next is not nil")
		}

		if rThird.prev != rSecond {
			t.Errorf("rThird prev is not rSecond")
		}
	})

	t.Run("Empty list", func(t *testing.T) {
		l := createList([]int{})
		reversed := ReverseList(l)

		if reversed.Count() != 0 {
			t.Fatalf("Unexpected changes: nodes are added")
		}

		if reversed.head != nil {
			t.Errorf("Empty list has head")
		}

		if reversed.tail != nil {
			t.Errorf("Empty list has tail")
		}
	})

	t.Run("Long list", func(t *testing.T) {
		l := createList([]int{1, 2, 3, 4, 5, 6})
		reversed := ReverseList(l)

		rHead := reversed.head
		if rHead.value != 6 {
			t.Errorf("rHead value is incorrect: %d, expected: %d", rHead.value, 6)
		}

		rTail := reversed.tail
		if rTail.value != 1 {
			t.Errorf("rTail value is incorrect: %d, expected: %d", rTail.value, 1)
		}
	})
}

func TestHasCycle(t *testing.T) {
	t.Run("Has no cycles", func(t *testing.T) {
		l := createList([]int{1, 2, 3})
		
		has:= HasCycle(l)

		if has {
			t.Errorf("The list does not definitely have cylces")
		}
	})

	t.Run("Has backward cycle", func(t *testing.T) {
		l := createList([]int{1, 2, 3, 4, 5})
		l.head.next.prev = l.tail // 2.prev == 5

		has := HasCycle(l)

		if !has {
			t.Errorf("The list does have backward cycle")
		}
	})

	t.Run("Has forward cycle", func(t *testing.T) {
		l := createList([]int{1, 2, 3, 4, 5})
		l.tail.prev.next = l.head.next // 4.next == 2

		has := HasCycle(l)

		if !has {
			t.Errorf("The list does have forward cycle")
		}
	})

	t.Run("Has multiple cylces", func(t *testing.T) {
		l := createList([]int{1, 2, 3, 4, 5})
		l.head.prev = l.tail // 1.prev == 5
		l.tail.prev.next = l.head.next // 4.next == 2

		has := HasCycle(l)

		if !has {
			t.Errorf("The list even has multiple cycles")
		}
	})

	t.Run("Empty list", func(t *testing.T) {
		l := createList([]int{})
		
		has := HasCycle(l)

		if has {
			t.Errorf("Empty list cannot have cycles")
		}
	})

	t.Run("One-element list with forward cycle", func(t *testing.T) {
		l := createList([]int{1})
		l.head.next = l.head

		has := HasCycle(l)

		if !has {
			t.Errorf("One-element list can be cycled to itself (forward)")
		}
	})

	t.Run("One-element list with backward cycle", func(t *testing.T) {
		l := createList([]int{1})
		l.head.prev = l.head

		has := HasCycle(l)

		if !has {
			t.Errorf("One-element list can be cycled to itself (backward)")
		}
	})

	t.Run("One-element list without cycle", func(t *testing.T) {
		l := createList([]int{1})

		has := HasCycle(l)

		if has {
			t.Errorf("The list has no cycles")
		}
	})
}

func TestSortList(t *testing.T) {
	t.Run("Unsorted list", func(t *testing.T) {
		l := createList([]int{3, 2, 1})
		sorted	 := SortList(l)

		sFirst := sorted.head
		sSecond := sFirst.next
		sThird := sSecond.next

		if sFirst.value != 1 || sSecond.value != 2 || sThird.value != 3 {
			t.Errorf("List is not sorted: %d, %d, %d", sFirst.value, sSecond.value, sThird.value)
		}
	})

	t.Run("Sorted list", func(t *testing.T) {
		l := createList([]int{1, 2, 3})
		sorted := SortList(l)

		sFirst := sorted.head
		sSecond := sFirst.next
		sThird := sSecond.next

		if sFirst.value != 1 || sSecond.value != 2 || sThird.value != 3 {
			t.Errorf("List is unsorted now: %d, %d, %d", sFirst.value, sSecond.value, sThird.value)
		}
	})
}

func TestMergeLists(t *testing.T) {
	t.Run("One-element lists of the same values", func(t *testing.T) {
		first := createList([]int{1})
		second :=createList([]int{1})

		result := MergeLists(first, second) // 1, 1

		rFirst := result.head
		rSecond := rFirst.next

		if result.Count() != 2 {
			t.Errorf("Merged list has unexpected count of elements: %d", result.Count())
		}

		if rFirst.value != 1 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 1)
		}

		if rFirst.prev != nil {
			t.Errorf("rHead has prev")
		}
	
		if rSecond.value != 1 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 1)
		}

		if rSecond.next != nil {
			t.Errorf("rTail has next")
		}
	
		if rSecond.prev != rFirst {
			t.Errorf("rSecond has incorrect prev")
		}
	})

	t.Run("One-element lists of different values", func(t *testing.T) {
		first := createList([]int{7})
		second :=createList([]int{3})

		result := MergeLists(first, second) // 3, 7

		rFirst := result.head
		rSecond := rFirst.next

		if rFirst.value != 3 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 3)
		}

		if rFirst.prev != nil {
			t.Errorf("rHead has prev")
		}
	
		if rSecond.value != 7 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 7)
		}

		if rSecond.next != nil {
			t.Errorf("rTail has next")
		}
	
		if rSecond.prev != rFirst {
			t.Errorf("rSecond has incorrect prev")
		}
	})

	t.Run("First list is empty", func(t *testing.T) {
		first := createList([]int{})
		second := createList([]int{3, 2})

		result := MergeLists(first, second)

		rFirst := result.head
		rSecond := rFirst.next

		if rFirst.value != 2 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 2)
		}

		if rFirst.prev != nil {
			t.Errorf("rHead has prev")
		}
	
		if rSecond.value != 3 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 3)
		}

		if rSecond.next != nil {
			t.Errorf("rTail has next")
		}
	
		if rSecond.prev != rFirst {
			t.Errorf("rSecond has incorrect prev")
		}
	})

	t.Run("Second list is empty", func(t *testing.T) {
		first := createList([]int{3, 2})
		second := createList([]int{})

		result := MergeLists(first, second)

		rFirst := result.head
		rSecond := rFirst.next

		if rFirst.value != 2 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 2)
		}

		if rFirst.prev != nil {
			t.Errorf("rHead has prev")
		}
	
		if rSecond.value != 3 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 3)
		}

		if rSecond.next != nil {
			t.Errorf("rTail has next")
		}
	
		if rSecond.prev != rFirst {
			t.Errorf("rSecond has incorrect prev")
		}
	})

	t.Run("Both lists are empty", func(t *testing.T) {
		first := createList([]int{})
		second := createList([]int{})

		result := MergeLists(first, second)

		rHead := result.head
		rTail := result.tail

		if rHead != nil {
			t.Errorf("rHead must be nil")
		}

		if rTail != nil {
			t.Errorf("rTail must be nil")
		}
	})

	t.Run("Unsorted lists of different length", func(t *testing.T) {
		first := createList([]int{2, 4, 0, 11}) // 0, 2, 4, 11
		second := createList([]int{5, 3, 2}) // 2, 3, 5

		result := MergeLists(first, second) // 0, 2, 2, 3, 4, 5, 11

		rFirst := result.head
		rSecond := rFirst.next
		rThird := rSecond.next
		rFourth := rThird.next
		rFifth := rFourth.next
		rSixth := rFifth.next
		rSeventh := rSixth.next
		
		if rFirst.value != 0 {
			t.Errorf("rFirst value is incorrect: %d, expected: %d", rFirst.value, 0)
		}

		if rFirst.prev != nil {
			t.Errorf("rHead has prev")
		}
	
		if rSecond.value != 2 {
			t.Errorf("rSecond value is incorrect: %d, expected: %d", rSecond.value, 2)
		}

		if rSecond.prev != rFirst {
			t.Errorf("rSecond has incorrect prev")
		}
	
		if rThird.value != 2 {
			t.Errorf("rThird value is incorrect: %d, expected: %d", rThird.value, 2)
		}

		if rThird.prev != rSecond {
			t.Errorf("rThird has incorrect prev")
		}
	
		if rFourth.value != 3 {
			t.Errorf("rFourth value is incorrect: %d, expected: %d", rFourth.value, 3)
		}

		if rFourth.prev != rThird {
			t.Errorf("rFourth has incorrect prev")
		}
	
		if rFifth.value != 4 {
			t.Errorf("rFifth value is incorrect: %d, expected: %d", rFifth.value, 4)
		}

		if rFifth.prev != rFourth {
			t.Errorf("rFifth has incorrect prev")
		}
	
		if rSixth.value != 5 {
			t.Errorf("rSixth value is incorrect: %d, expected: %d", rSixth.value, 5)
		}

		if rSixth.prev != rFifth {
			t.Errorf("rSixth has incorrect prev")
		}
	
		if rSeventh.value != 11 {
			t.Errorf("rSeventh value is incorrect: %d, expected: %d", rSeventh.value, 11)
		}

		if rSeventh.next != nil {
			t.Errorf("rTail has next")
		}
	
		if rSeventh.prev != rSixth {
			t.Errorf("rSeventh has incorrect prev")
		}
	})
}


