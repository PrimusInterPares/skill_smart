package main

import (
	"testing"
)

func createList(values []int) *LinkedList {
	list := &LinkedList{}
	for _, value := range values {
		list.AddInTail(Node{value: value})
	}

	return list
}

func TestLinkedList_AddInTail(t *testing.T) {
	testCases := []struct {
		title string
		initialValues []int
		newNode Node
		expectedHeadValue int
		expectedTailValue int
	}{
		{
			title: "Empty list",
			initialValues: []int{},
			newNode: Node{value: 4},
			expectedHeadValue: 4,
			expectedTailValue: 4,
		},
		{
			title: "List with elements",
		  initialValues: []int{1, 2, 3},
			newNode: Node{value: 4},
			expectedHeadValue: 1,
			expectedTailValue: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			list := createList(tc.initialValues)
			
			list.AddInTail(tc.newNode)

			headValue := list.head.value
			if headValue != tc.expectedHeadValue {
				t.Fatalf("Head value must be %d but got %d", headValue, tc.expectedHeadValue)
			}

			tailValue := list.tail.value
			if tailValue != tc.expectedTailValue {
				t.Errorf("Tail value must be %d but got %d", tc.expectedTailValue, tailValue)
			}
		})
	}
}

func TestLinkedList_Count(t *testing.T) {
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

func TestLinkedList_Find(t *testing.T) {
	list := createList([]int{1, 2, 3})

	t.Run("Found", func(t *testing.T) {
		findValue := 2
		resultNode, resultErr := list.Find(findValue)

		if resultErr != nil {
			t.Fatalf("Unexpected error: %v", resultErr)
		}

		if resultNode.value != findValue {
			t.Errorf("Wrong value found: %d, expected: %d", resultNode.value, findValue)
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

func TestLinkedList_FindAll(t *testing.T) {
	list := createList([]int{0, 1, 1, 2, 3, 5})

	t.Run("Array found", func(t *testing.T) {
		findValue := 1
		result := list.FindAll(findValue)

		if len(result) != 2 {
			t.Fatalf("More or less elements found. Result: %v, expected: [1, 1]", result)
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

func TestLinkedList_Delete(t *testing.T) {
	testCases := []struct {
		title string
		initialValues []int
		toDelete int
		all bool
		expectedCount int
		expectedHead int
		expectedTail int
	}{
		{
			title: "Delete head",
			initialValues: []int{1, 3},
			toDelete: 1,
			all: false,
			expectedCount: 1,
			expectedHead: 3,
			expectedTail: 3,
		},
		{
			title: "Delete tail",
			initialValues: []int{1, 3},
			toDelete: 3,
			all: false,
			expectedCount: 1,
			expectedHead: 1,
			expectedTail: 1,
		},
			{
			title: "Delete from the middle",
			initialValues: []int{1, 2, 3},
			toDelete: 2,
			all: false,
			expectedCount: 2,
			expectedHead: 1,
			expectedTail: 3, 
		},
		{
			title: "Node is not found",
			initialValues: []int{1, 2, 3},
			toDelete: 5,
			all: false,
			expectedCount: 3,
			expectedHead: 1, 
			expectedTail: 3,
		},
		{
			title: "Delete all inclusions",
			initialValues: []int{1, 2, 1, 3, 1, 4},
			toDelete: 1,
			all: true,
			expectedCount: 3,
			expectedHead: 2,
			expectedTail: 4,
		},
		{
			title: "Multiple heads to delete",
			initialValues: []int{1, 1, 1, 1, 2},
			toDelete: 1,
			all: true,
			expectedCount: 1,
			expectedHead: 2,
			expectedTail: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			list := createList(tc.initialValues)
			list.Delete(tc.toDelete, tc.all)

			if list.Count() != tc.expectedCount {
				t.Errorf("Got %d elements after deletion, expected %d", list.Count(), tc.expectedCount)
			}

			if list.head.value != tc.expectedHead {
				t.Errorf("Expected head value: %d, got: %d", tc.expectedHead, list.head.value)
			}

			if list.tail.value != tc.expectedTail {
				t.Errorf("Expected tail value: %d, got: %d", tc.expectedTail, list.tail.value)
			}
		})
	}

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
	})
	
	t.Run("One-element list", func(t *testing.T) {
		list := createList([]int{1})
		list.Delete(1, false)

		if list.Count() != 0 {
			t.Fatalf("Not all elements were deleted, %d elements left.", list.Count())
		}
	})
}

func TestLinkedList_InsertFirst(t *testing.T) {
	t.Run("Regular", func(t *testing.T) {
		list := createList([]int{2, 3})
		
		insertValue := 4
		list.InsertFirst(Node{value: insertValue})

		if list.Count() != 3 {
			t.Errorf("Incorrect elements count: %d, expected: %d", list.Count(), insertValue)
		}

		if list.head.value != insertValue {
			t.Errorf("Incorrect head value: %d, expected: %d", list.head.value, insertValue)
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


