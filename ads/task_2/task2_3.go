package main

import(
	"testing"
)

func createList(values []int) *LinkedList2 {
	list := &LinkedList2{}
	for _, value := range values {
		list.AddInTail(Node{value: value})
	}

	return list
}

func TestLinkedList2_AddInTail(t *testing.T) {
	newNode := Node{value: 4}

	t.Run("Empty list", func(t *testing.T) {
		list := createList([]int{})
		expectedHeadValue := 4
		expectedTailValue := 4

		list.AddInTail(newNode)
		
		headValue := list.head.value
		if headValue != expectedHeadValue {
			t.Fatalf("Head value must be %d but got %d", expectedHeadValue, headValue)
		}
		
		tailValue := list.tail.value
		if tailValue != expectedTailValue {
			t.Errorf("Tail value must be %d but got %d", expectedTailValue, tailValue)
		}
		
		prev := list.tail.prev
		if prev != nil {
			t.Errorf("Prev must be absent but got %d", prev.value)
		}
	})
	
	t.Run("List with elements", func(t *testing.T) {
		list := createList([]int{1, 2, 3})
		expectedHeadValue := 1
		expectedTailValue := 4
		expectedPrevValue := 3

		list.AddInTail(newNode)
		
		headValue := list.head.value
		if headValue != expectedHeadValue {
			t.Fatalf("Head value must be %d but got %d", expectedHeadValue, headValue)
		}
		
		tailValue := list.tail.value
		if tailValue != expectedTailValue {
			t.Errorf("Tail value must be %d but got %d", expectedTailValue, tailValue)
		}
		
		prev := list.tail.prev
		if prev == nil {
			t.Fatalf("Prev must be present but got `nil`")
		}

		prevValue := prev.value
		if prevValue != expectedPrevValue {
			t.Errorf("Prev value must be %d but got %d", expectedPrevValue, prevValue)
		}
	})
}

func TestLinkedList2_Count(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		list := createList([]int{}) 
		result := list.Count()
		expectedResult := 0

		if result != expectedResult {
			t.Errorf("Count for the empty list = %d, expected: %d", result, expectedResult)
		}
	})

	t.Run("List with multiple values", func(t *testing.T) {
		list := createList([]int{1, 2, 3})
		result := list.Count()
		expectedResult := 3

		if result != expectedResult {
			t.Errorf("Count for the list = %d, expected: %d", result, expectedResult)
		}
	})
}

func TestLinkedList2_Find(t *testing.T) {
	list := createList([]int{1, 2, 3, 2})

	t.Run("Found", func(t *testing.T) {
		findValue := 2
		resultNode, resultErr := list.Find(findValue)

		if resultErr != nil {
			t.Fatalf("Unexpected error: %v", resultErr)
		}

		if resultNode.value != findValue {
			t.Errorf("Wrong value found: %d, expected: %d", resultNode.value, findValue)
		}

		if &resultNode == list.tail {
			t.Fatalf("The second element was found - first must be")
		}
	})

	t.Run("Not found", func(t *testing.T) {
		resultNode, resultErr := list.Find(5)

		if resultErr == nil {
			t.Fatalf("Value `%d` found, error expected", resultNode.value)
		}

		if resultErr.Error() != "not found" {
			t.Errorf("Unexpected error message: `%v`", resultErr.Error())
		}
	})
}

func TestLinkedList2_FindAll(t *testing.T) {
	list := createList([]int{1, 0, 1, 2, 3, 5, 1})

	t.Run("Array found", func(t *testing.T) {
		findValue := 1
		result := list.FindAll(findValue)

		if len(result) != 3 {
			t.Fatalf("More or less elements found. Result: %v, expected: [1, 1, 1]", result)
		}
	})

	t.Run("None found", func(t *testing.T) {
		findValue := 100
		result := list.FindAll(findValue)

		if len(result) != 0 {
			t.Errorf("Expected empty array, got: %v", result)
		}
	})
}

func TestLinkedList2_Delete(t *testing.T) {
	t.Run("Delete head", func(t *testing.T) {
		expectedCount := 1
		expectedHead := 3
		expectedTail := 3
		list := createList([]int{1, 3})
		list.Delete(1, false)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}

		if list.tail.prev != nil {
			t.Errorf("Tail prev is not nil")
		}

		if list.head.next != nil {
			t.Errorf("Head next is not nil")
		}
	})
	
	t.Run("Delete tail", func(t *testing.T) {
		expectedCount := 1
		expectedHead := 1
		expectedTail := 1
		list := createList([]int{1, 3})
		list.Delete(3, false)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}

		if list.tail.prev != nil {
			t.Errorf("Tail prev is not nil")
		}

		if list.head.next != nil {
			t.Errorf("Head next is not nil")
		}
	})
	
	t.Run("Delete from the middle", func(t *testing.T) {
		expectedCount := 2
		expectedHead := 1
		expectedTail := 3
		list := createList([]int{1, 2, 3})
		list.Delete(2, false)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}

		if list.head.next.value != expectedTail {
			t.Errorf("Head next is incorrect: %d, expecrted: %d", list.head.next.value, expectedTail)
		}

		if list.tail.prev.value != expectedHead {
			t.Errorf("Tail prev is incorrect: %d, expecrted: %d", list.tail.prev.value, expectedHead)
		}
	})
	
	t.Run("Node is not found", func(t *testing.T) {
		expectedCount := 3
		expectedHead := 1 
		expectedTail := 3
		list := createList([]int{1, 2, 3})
		list.Delete(5, false)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}
	})
	
	t.Run("Delete all inclusions", func(t *testing.T) {
		expectedCount := 3
		expectedHead := 2
		expectedTail := 4
		list := createList([]int{1, 2, 1, 3, 1, 4})
		list.Delete(1, true)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}

		if list.head.prev != nil {
			t.Errorf("Head prev is not nil")
		}

		if list.tail.next != nil {
			t.Errorf("Tail next is not nil")
		}

		secondNode := list.head.next
		if secondNode.prev != list.head {
			t.Errorf("Middle node prev is not head")
		}

		if secondNode.next != list.tail {
			t.Errorf("Middle node next is not tail")
		}
	})
	
	t.Run("Multiple heads to delete", func(t *testing.T) {
		expectedCount := 1
		expectedHead := 2
		expectedTail := 2
		list := createList([]int{1, 1, 1, 1, 2})
		list.Delete(1, true)

		if list.Count() != expectedCount {
			t.Errorf("Got %d elements after deletion, expected %d", list.Count(), expectedCount)
		}

		if list.head.value != expectedHead {
			t.Errorf("Expected head value: %d, got: %d", expectedHead, list.head.value)
		}

		if list.tail.value != expectedTail {
			t.Errorf("Expected tail value: %d, got: %d", expectedTail, list.tail.value)
		}

		if list.head.next != nil {
			t.Errorf("Head next is not nil")
		}

		if list.tail.prev != nil {
			t.Errorf("Tail prev is not nil")
		}
	})

	t.Run("Delete from empty list", func(t *testing.T) {
		list := createList([]int{})
		list.Delete(5, false)

		if list.Count() > 0 {
			t.Fatalf("Deletion caused insert!")
		}
	})

	t.Run("Delete the whole list", func(t *testing.T) {
		list := createList([]int{1, 1, 1})
		list.Delete(1, true)

		if list.Count() != 0 {
			t.Fatalf("Not all elements were deleted, %d elements left.", list.Count())
		}

		if list.head != nil || list.tail != nil {
			t.Errorf("List is not empty")
		}
	})
	
	t.Run("One-element list", func(t *testing.T) {
		list := createList([]int{1})
		list.Delete(1, false)

		if list.Count() != 0 {
			t.Fatalf("Not all elements were deleted, %d elements left.", list.Count())
		}
	})
}

func TestLinkedList2_Insert(t *testing.T) {
	list := createList([]int{1, 2, 4, 5})
	add := Node{value: 3}

	t.Run("In the middle", func(t *testing.T) {
		after := list.head.next

		list.Insert(after, add)

		insertedNode := after.next
		if insertedNode.value != add.value {
			t.Errorf("Node was not iserted properly after %d. Current next: %d, must be: %d", 
				after.value, insertedNode.value, add.value)
		}

		if insertedNode.prev != after {
			t.Errorf("Prev of the inserted node is incorrect. Current value: %d, expected: %d",
				insertedNode.prev.value, after.value)
		}
	})

	t.Run("One-element list", func(t *testing.T) {
		list := createList([]int{1})
		after := list.head

		list.Insert(after, add)

		if list.tail.value != add.value {
			t.Errorf("Node was not inserted after head. Current next: %d, must be: %d", list.tail.value, add.value)
		}

		if list.tail.prev == nil {
			t.Fatalf("Inserted node must have prev")
		}

		if list.tail.prev != after {
			t.Errorf("Prev of the inserted node is incorrect. Current value: %d, expected: %d",
				list.tail.prev.value, after.value)
		}
	})
}

func TestLinkedList2_InsertFirst(t *testing.T) {
	t.Run("Regular", func(t *testing.T) {
		list := createList([]int{2, 3})
		
		insertValue := 4
		initialHead := list.head
		list.InsertFirst(Node{value: insertValue})

		if list.Count() != 3 {
			t.Errorf("Incorrect elements count: %d, expected: %d", list.Count(), 3)
		}

		if list.head.value != insertValue {
			t.Errorf("Incorrect head value: %d, expected: %d", list.head.value, insertValue)
		}

		if list.head == initialHead {
			t.Fatalf("Head is not changed")
		}

		if initialHead.prev.value != insertValue {
			t.Errorf("Incorrect prev value: %d, expected: %d", initialHead.prev.value, insertValue)
		}
	})

	t.Run("Insert in empty", func(t *testing.T) {
		list := createList([]int{})
		
		list.InsertFirst(Node{value: 1})
		
		if list.Count() == 0 {
			t.Errorf("Nothing was inserted.")
		}

		if list.head.value != 1 {
			t.Errorf("Incorrect head value: %d, expected: %d", list.head.value, 1)
		}

		if list.tail.value != 1 {
			t.Errorf("Incorrect tail value: %d, expected: %d", list.tail.value, 1)
		}
	})
}

