package main

import(
	"testing"
)

func Test_SumLists(t *testing.T) {
	t.Run("Lists have different sizes", func(t *testing.T) {
		first := createList([]int{1, 2})
		second := createList([]int{1, 2, 3})

		result, err := SumLists(first, second)
		
		if err == nil {
			t.Errorf("Error must be present!")
		}
		
		if result.Count() != 0 {
			t.Errorf("Returned list must be empty")
		}
	})

	t.Run("Lists are summed", func(t *testing.T) {
		first := createList([]int{1, 2, 3})
		second := createList([]int{4, 5, 6})

		result, err := SumLists(first, second)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err.Error())
		}

		if result.Count() != 3 {
			t.Errorf("Incorrect output list size: %d", result.Count())
		}


		if result.head.value != 5 || result.head.next.value != 7 || result.head.next.next.value != 9 {
			t.Errorf("Incorrect calculation.")
		}
	})
}


