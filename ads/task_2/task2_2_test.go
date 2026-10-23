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

}

func TestSortList(t *testing.T) {

}

func TestMergeLists(t *testing.T) {

}


