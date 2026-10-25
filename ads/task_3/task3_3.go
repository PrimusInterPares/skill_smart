package main

import(
	"testing"
)

func createDA[T any](items []T) *DynArray[T] {
	da := &DynArray[T]{}

	if len(items) == 0 {
		da.Init()

		return da
	}

	initCapacity := max(minimalBufferSize, len(items))
	da.MakeArray(initCapacity)
	copy(da.array, items)
	da.count = len(items)

	return da
}

func makeRange(minVal, maxVal int) []int {
	r := make([]int, maxVal - minVal + 1)
	
	for i := range r {
		r[i] = minVal + i
	}

	return r
}

func TestDynArray_Insert(t *testing.T) {
	t.Run("Buffer remains unchanged", func(t *testing.T) {
		arr := createDA(makeRange(2, 6))

		err:= arr.Insert(100, 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
		arr := createDA(makeRange(1, 16))

		err := arr.Insert(100, 12)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
		arr := createDA([]int{1, 2, 3})

		err := arr.Insert(100, -1)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '-1'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index > count", func(t *testing.T) {
		arr := createDA([]int{1, 2, 3})

		err := arr.Insert(100, 4)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '4'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index == count", func(t *testing.T) {
		arr := createDA([]int{1, 2, 3})

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

func TestDynArray_Remove(t *testing.T) {
	t.Run("Buffer remains unchanged", func(t *testing.T) {
		arr := createDA([]int{1, 2, 3})

		err := arr.Remove(1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
		arr := createDA(makeRange(0, 31))
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
		arr := createDA(makeRange(1, 8))
		
		err := arr.Remove(1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if arr.capacity < minimalBufferSize {
			t.Errorf("capacity must be at least %d, current: %d", minimalBufferSize, arr.capacity)
		}

		if arr.count != 7 {
			t. Errorf("expected count: %d, current: %d", 7, arr.count)
		}
	})

	t.Run("Negative index", func(t *testing.T) {
		arr := createDA([]int{1, 2, 3})

		err := arr.Remove(-1)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '-1'" {
			t.Errorf("unexpected error")
		}
	})

	t.Run("Index > count", func(t *testing.T) {
		arr := createDA([]int{1, 2, 3})

		err := arr.Remove(4)

		if err == nil {
			t.Errorf("must be bad index error")
		}

		if err.Error() != "bad index '4'" {
			t.Errorf("unexpected error")
		}
	})
}


