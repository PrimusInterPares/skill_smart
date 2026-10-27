package main

import(
	"testing"
)

func createBDA[T any](items []T) *BankDynArray[T] {
	bda := NewBankDynArray[T]()

	for _, item := range items {
		bda.Append(item)
	}

	return bda
}

func TestBankDynArray_Insert(t *testing.T) {
	t.Run("Buffer remains unchanged", func(t *testing.T) {
		arr := createBDA(makeRange(2, 6))

		arr.Insert(100, 1)

		if v, _ := arr.GetItem(0); v != 2 {
			t.Errorf("arr[0] is expected to eq 2, actual: %d", v)
		}

		if v, _ := arr.GetItem(1); v != 100 {
			t.Errorf("arr[1] is expected to eq 100, actual: %d", v)
		}

		if v, _ := arr.GetItem(2); v != 3 {
			t.Errorf("arr[2] is expected to be 3, actual: %d", v)
		}

		if arr.capacity != 16 {
			t.Errorf("capacity must not have been changed, current: %d", arr.capacity)
		}

		if arr.count != 6 {
			t. Errorf("expected count: %d", 6)
		}
	})

	t.Run("Buffer changes", func(t *testing.T) {
		arr := createBDA(makeRange(1, 16))

		arr.Insert(100, 12)

		if v, _ := arr.GetItem(11); v != 12 {
			t.Errorf("arr[11] is expected to eq 12, actual: %d", v)
		}

		if v, _ := arr.GetItem(12); v != 100 {
			t.Errorf("arr[12] is expected to eq 100, actual: %d", v)
		}

		if v, _ := arr.GetItem(13); v != 13 {
			t.Errorf("arr[13] is expected to be 13, actual: %d", v)
		}

		if arr.capacity == 16 {
			t.Errorf("capacity must have been changed to %d, current: %d", 16 * increaseBufferRatio, arr.capacity)
		}

		if arr.count != 17 {
			t. Errorf("expected count: %d, current: %d", 17, arr.count)
		}
	})

	t.Run("Negative index", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		err := arr.Insert(100, -1)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '-1'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index > count", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		err := arr.Insert(100, 4)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '4'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index == count", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		err := arr.Insert(100, 3)

		if v, _ := arr.GetItem(3); v != 100 {
			t.Errorf("value 100 must be added to position 3")
		}

		if arr.count != 4 {
			t.Errorf("count must eq %d, current: %d", 4, arr.count)
		}

		if err != nil {
			t.Errorf("unexpected error")
		}
	})
}

func TestBankDynArray_Remove(t *testing.T) {
	t.Run("Buffer remains unchanged", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		arr.Remove(1)

		if v, _ := arr.GetItem(0); v != 1 {
			t.Errorf("arr[0] is expected to eq 1, actual: %d", v)
		}

		if v, _ := arr.GetItem(1); v != 3 {
			t.Errorf("arr[1] is expected to eq 3, actual: %d", v)
		}

		if _, err := arr.GetItem(2); err == nil {
			t.Errorf("must be bad index error")
		}

		if arr.capacity != 16 {
			t.Errorf("capacity must not have been changed, current: %d", arr.capacity)
		}

		if arr.count != 2 {
			t. Errorf("expected count: %d, current: %d", 2, arr.count)
		}
	})

	t.Run("Buffer changes", func(t *testing.T) {
		arr := createBDA(makeRange(0, 31))
		oldCount := arr.count
		oldCapacity := arr.capacity

		for i := arr.count - 1; 14 < i; i-- {
			if err := arr.Remove(i); err != nil {
				t.Fatalf("unexpected error on iteration %d: %v", i, err)
			}
		}

		expectedCount := 32 - 17 // 15
		if arr.count != expectedCount {
			t. Errorf("expected count: %d, current: %d, initial: %d", expectedCount, arr.count, oldCount)
		}
		
		// 15 * 2 < 32 -> decreaseBuffer
		if arr.capacity >= oldCapacity {
			t.Errorf("capacity must have been changed from %d", oldCapacity)
		}
	})

	t.Run("Buffer cannot be less than minimalBufferSize", func(t *testing.T) {
		arr := createBDA(makeRange(1, 8))
		
		arr.Remove(1)

		if arr.capacity < minimalBufferSize {
			t.Errorf("capacity must be at least %d, current: %d", minimalBufferSize, arr.capacity)
		}

		if arr.count != 7 {
			t. Errorf("expected count: %d, current: %d", 7, arr.count)
		}
	})

	t.Run("Negative index", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		err := arr.Remove(-1)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '-1'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index > count", func(t *testing.T) {
		arr := createBDA([]int{1, 2, 3})

		err := arr.Remove(4)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '4'" {
			t.Errorf("unexpected error")
		}
	})
}

//--------------------------------------------------------------------------------------

func TestMultiDynArray_NewMultiDynArray(t *testing.T) {
	t.Run("Zero depth", func(t *testing.T) {
		_, err := NewMultiDynArray[int](0)

		if err == nil { t.Fatalf("error is expected") }
		if err.Error() != "depth must be positive, got: 0" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	
	t.Run("Negative depth", func(t *testing.T) {
		_, err := NewMultiDynArray[int](-1)

		if err == nil { t.Fatalf("error is expected") }

		if err.Error() != "depth must be positive, got: -1"{
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestMultiDynArray_GetItem(t *testing.T) {
	arr, _ := NewMultiDynArray[int](3) // 2x3x4

	t.Run("0 indices provided", func(t *testing.T) {
		_, err := arr.GetItem()

		if err == nil { t.Fatalf("error is expected") }
		if err.Error() != "wrong number of indices: 0, expected: 3" { t.Errorf("unexpected error: %v", err) }
	})

	t.Run("Less indices than needed", func(t *testing.T) {
		_, err := arr.GetItem(1)

		if err == nil { t.Fatalf("error expected") }
		if err.Error() != "wrong number of indices: 1, expected: 3" { t.Errorf("unexpected error: %v", err) }
	})

	t.Run("Too many indices", func(t *testing.T) {
		_, err := arr.GetItem(1, 2, 3, 4)

		if err == nil { t.Fatalf("error is expected") }
		if err.Error() != "wrong number of indices: 4, expected: 3" { t.Errorf("unexpected error: %v", err) }
	})

	t.Run("Positive case", func(t *testing.T) {
		err := arr.Append(5, 0, 0)
		if err != nil { t.Fatalf("unexpected error during Append: %v", err) }

		value, err := arr.GetItem(0, 0, 0)
		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if value != 5 { t.Errorf("expected value: 5, got: %d", value) }
	})
}

func TestMultiDynArray_Append(t *testing.T) {
	t.Run("1D array", func(t *testing.T) {
    arr, err := NewMultiDynArray[int](1)
		if err != nil { t.Fatal(err) }

    err = arr.Append(42)
    if err != nil { t.Fatalf("unexpected error: %v", err) }

    val, err := arr.GetItem(0)
    if err != nil { t.Fatal(err) }
    if val != 42 { t.Errorf("expected 42, got %d", val) }
	})

	t.Run("Triggers buffer increase", func(t *testing.T) {
    arr, err := NewMultiDynArray[int](1)
    if err != nil { t.Fatal(err) }

    for i := range 20 {
      err := arr.Append(i)
      if err != nil { t.Fatalf("unexpected error on iteration %d: %v", i, err) }
    }

    for i := range 20 {
      val, err := arr.GetItem(i)
      if err != nil { t.Fatalf("unexpected error getting item at index [0, %d]: %v",  i, err) }
			if val != i { t.Errorf("expected %d at index [0, %d], got: %d", i, i, val) }
    }
	})

	t.Run("Wrong number of indices", func(t *testing.T) {
    arr, _ := NewMultiDynArray[int](3)

    err := arr.Append(5, 1)
    if err == nil { t.Error("expected error") }
		if err.Error() != "expected 2 indices, got 1" {
			t.Errorf("unexpected error: %v", err)
		}

    err = arr.Append(5, 1, 1, 1)
    if err == nil { t.Error("expected error") }
		if err.Error() != "expected 2 indices, got 3" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Out-of-bounds navigation indices", func(t *testing.T) {
    arr, _ := NewMultiDynArray[int](3)

    err := arr.Append(5, 0)
    if err == nil { t.Error("expected error") }
		if err.Error() != "expected 2 indices, got 1" { t.Errorf("unexpected error: %v", err) }
	})

	t.Run("With string type", func(t *testing.T) {
    arr, err := NewMultiDynArray[string](2)
    if err != nil { t.Fatal(err) }

    err = arr.Append("hello", 0)
    if err != nil { t.Fatal(err) }

    val, err := arr.GetItem(0, 0)
		if err != nil { t.Fatalf("unexpected error: %v",err) }
		if val != "hello" { t.Errorf("expected 'hello', got: '%s'", val) }
	})

	t.Run("Buffer increase", func(t *testing.T) {
    arr, _ := NewMultiDynArray[int](2)

    for i := range 20 {
	    err := arr.Append(i * 10, i)
      if err != nil { t.Fatalf("unexpected error on iteration %d: %v", i, err) }
    }

    if arr.nodes.count != 20 {
    	t.Errorf("expected nodes count 20, got %d", arr.nodes.count)
    }
    if arr.nodes.capacity < 20 {
    	t.Errorf("nodes capacity must have increased, got %d", arr.nodes.capacity)
    }

    for i := range 20 {
	    val, err := arr.GetItem(i, 0)
      if err != nil { t.Fatalf("unexpected error at [%d, 0]: %v", i, err) }
      if val != i * 10 { t.Errorf("expected %d at [%d, 0], got %d", i * 10, i, val) }
    }
})
}

func TestMultiDynArray_Insert(t *testing.T) {
  t.Run("Insert into existing leaf", func(t *testing.T) {
		arr, _ := NewMultiDynArray[int](2)

		arr.Append(10, 0)
		arr.Append(20, 0) // [10, 20]

    err := arr.Insert(99, 0, 1)
    if err != nil { t.Fatal(err) }

		// [10, 99, 20]
    val, _ := arr.GetItem(0, 0)
    if val != 10 { t.Errorf("expected 10, got %d", val) }
 
		val, _ = arr.GetItem(0, 1)
    if val != 99 { t.Errorf("expected 99, got %d", val) }
 
    val, _ = arr.GetItem(0, 2)
    if val != 20 { t.Errorf("expected 20, got %d", val) }
  })

  t.Run("Insert into non-existent node", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](2)
    arr.Append(10, 0) // только лист [0]

    err := arr.Insert(42, 5, 0)
    if err == nil { t.Error("expected error") }
  })

  t.Run("Insert at the end the of existing leaf", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](2)
    arr.Append(10, 0)

    err := arr.Insert(20, 0, 1)
    if err != nil { t.Fatal(err) }

		val, _ := arr.GetItem(0, 1)
    if val != 20 { t.Errorf("expected 20, got %d", val) }
    })
}

func TestMultiDynArray_Remove(t *testing.T) {
  t.Run("Remove from existing leaf", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](2)
		
    arr.Append(10, 0)
    arr.Append(20, 0)
    arr.Append(30, 0) // [10, 20, 30]

    err := arr.Remove(0, 1)
    if err != nil { t.Fatal(err) }

    val, _ := arr.GetItem(0, 0)
    if val != 10 { t.Errorf("expected 10, got %d", val) }

		val, _ = arr.GetItem(0, 1)
    if val != 30 { t.Errorf("expected 30, got %d", val) }
  })

  t.Run("Remove from non-existent node", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](2)
    arr.Append(10, 0)

    err := arr.Remove(5, 0)
    if err == nil { t.Error("expected error for non-existent node") }
  })

  t.Run("Wrong number of indices", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](3)
    arr.Append(10, 0, 0)

    err := arr.Remove(0, 0)
    if err == nil { t.Error("expected error for too few indices") }

    err = arr.Remove(0, 0, 0, 0)
    if err == nil { t.Error("expected error for too many indices") }
  })

	t.Run("Buffer decrease", func(t *testing.T) {
  	arr, _ := NewMultiDynArray[int](2)

    for i := range 40 {
        arr.Append(i, 0)
    }

    leaf, _ := arr.nodes.GetItem(0)
    capacityBefore := leaf.values.capacity

    for i := range 30 {
      err := arr.Remove(0, 0)
      if err != nil { t.Fatalf("unexpected error on iteration %d: %v", i, err) }
    }

    if leaf.values.count != 10 { t.Errorf("expected count 10, got %d", leaf.values.count) }

    if leaf.values.capacity >= capacityBefore {
        t.Errorf("capacity must have decreased from %d, got %d", capacityBefore, leaf.values.capacity)
    }

    for i := range 10 {
      val, err := arr.GetItem(0, i)
      if err != nil { t.Fatalf("unexpected error at [0, %d]: %v", i, err) }

			expected := 30 + i
      if val != expected { t.Errorf("expected %d at [0, %d], got %d", expected, i, val)}
    }
	})
}


