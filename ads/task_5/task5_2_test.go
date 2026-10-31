package main

import (
	"errors"
	"slices"
	"testing"
)

func TestQueue_Rotate(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		n        int
		expected []int
	}{
		{"just rotate", []int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{"rotate for 0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"n > size", []int{1, 2, 3}, 5, []int{3, 1, 2}},           // 5 % 3 = 2
		{"negative n", []int{1, 2, 3, 4}, -13, []int{4, 1, 2, 3}}, // -13 % 4 = -1
		{"empty queue", []int{}, 2, []int{}},
		{"one element", []int{1}, 10, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := &Queue[int]{items: tt.initial}

			queue.Rotate(tt.n)

			if !slices.Equal(queue.items, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, queue.items)
			}
		})
	}
}

func TestStackedQueue_Size(t *testing.T) {
	tests := []struct {
		name     string
		elems    []int
		expected int
	}{
		{"empty queue", []int{}, 0},
		{"one element", []int{5}, 1},
		{"queue with elements", []int{1, 2, 3, 4}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewStackedQueue[int]()

			for _, v := range tt.elems {
				queue.Enqueue(v)
			}

			qSize := queue.Size()

			if qSize != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, qSize)
			}
		})
	}

	t.Run("stackIn is empty", func(t *testing.T) {
		queue := NewStackedQueue[string]()

		queue.Enqueue("first")
		queue.Enqueue("second")
		_, _ = queue.Peek()

		if queue.stackIn.Size() != 0 {
			t.Fatal("stackIn is not empty")
		}
		if queue.Size() != 2 {
			t.Errorf("Expected 2, got %d", queue.Size())
		}
	})
}

func TestStackedQueue_Push(t *testing.T) {

}

func TestStackedQueue_Peek(t *testing.T) {
	queue := NewStackedQueue[int]()

	t.Run("empty queue", func(t *testing.T) {
		result, err := queue.Peek()

		if err == nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		var zero int
		if result != zero {
			t.Errorf("Result must be empty, got %v", result)
		}
		if !errors.Is(err, ErrEmptyQueue) {
			t.Errorf("Expected ErrEmptyQueue, got %v", err)
		}
	})

	t.Run("queue with elements", func(t *testing.T) {
		expectedHead := 10
		queue.Enqueue(expectedHead)
		queue.Enqueue(102030)
		queue.Enqueue(888)

		result, err := queue.Peek()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if queue.stackIn.Size() != 0 {
			t.Error("stackIn must be empty")
		}
		if result != expectedHead {
			t.Errorf("Expected %d, got %d", expectedHead, result)
		}
	})
}

func TestStackedQueue_Pop(t *testing.T) {
	queue := NewStackedQueue[int]()

	t.Run("empty queue", func(t *testing.T) {
		result, err := queue.Dequeue()

		if err == nil {
			t.Fatal("Error is expected")
		}
		if !errors.Is(err, ErrEmptyQueue) {
			t.Errorf("Expected ErrEmptyQueue, got %v", err)
		}
		var zero int
		if result != zero {
			t.Errorf("Result must be empty, got %d", result)
		}
	})

	t.Run("queue with elements", func(t *testing.T) {
		expectedHead := 999
		queue.Enqueue(expectedHead)
		queue.Enqueue(0)
		queue.Enqueue(-100)

		result, err := queue.Dequeue()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != expectedHead {
			t.Errorf("Expected %d, got %d", expectedHead, result)
		}
	})
}

func TestQueue_Reverse(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected []int
	}{
		{"simple queue", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"empty queue", []int{}, []int{}},
		{"one element", []int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue[int]{items: tt.initial}

			q.Reverse()

			if !slices.Equal(q.items, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, q.items)
			}
		})
	}
}

func TestRingQueue_Enqueue(t *testing.T) {
	t.Run("enqueue single element", func(t *testing.T) {
		queue := NewRingQueue[int](3)

		err := queue.Enqueue(10)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		val, _ := queue.Dequeue()
		if val != 10 {
			t.Errorf("Expected 10, got %d", val)
		}
	})

	t.Run("enqueue until full", func(t *testing.T) {
		queue := NewRingQueue[int](2)

		if err := queue.Enqueue(1); err != nil {
			t.Errorf("unexpected error on 1st enqueue: %v", err)
		}
		if err := queue.Enqueue(2); err != nil {
			t.Errorf("unexpected error on 2nd enqueue: %v", err)
		}
	})

	t.Run("enqueue when full returns ErrFullQueue", func(t *testing.T) {
		queue := NewRingQueue[int](1)

		_ = queue.Enqueue(1)

		err := queue.Enqueue(2)
		if !errors.Is(err, ErrFullQueue) {
			t.Errorf("expected ErrFullQueue, got %v", err)
		}
	})

	t.Run("enqueue after dequeue (wrap-around)", func(t *testing.T) {
		queue := NewRingQueue[int](3)
		_ = queue.Enqueue(1)
		_ = queue.Enqueue(2)

		_, _ = queue.Dequeue() // head moves to index 1
		_ = queue.Enqueue(3)   // index must be 0

		// FIFO: 2, 3
		v1, _ := queue.Dequeue()
		v2, _ := queue.Dequeue()
		if v1 != 2 || v2 != 3 {
			t.Errorf("Expected 2 and 3, got %d and %d", v1, v2)
		}
	})
}

func TestRingQueue_Dequeue(t *testing.T) {
	t.Run("dequeue from empty returns ErrEmptyQueue", func(t *testing.T) {
		queue := NewRingQueue[int](3)

		_, err := queue.Dequeue()
		if !errors.Is(err, ErrEmptyQueue) {
			t.Errorf("Expected ErrEmptyQueue, got %v", err)
		}
	})

	t.Run("dequeue single element", func(t *testing.T) {
		queue := NewRingQueue[int](3)

		_ = queue.Enqueue(42)

		val, err := queue.Dequeue()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("dequeue maintains FIFO order", func(t *testing.T) {
		queue := NewRingQueue[int](5)

		for i := 1; i <= 3; i++ {
			_ = queue.Enqueue(i)
		}

		for i := 1; i <= 3; i++ {
			val, err := queue.Dequeue()
			if err != nil {
				t.Errorf("Unexpected error on dequeue %d: %v", i, err)
			}
			if val != i {
				t.Errorf("Expected %d, got %d", i, val)
			}
		}
	})

	t.Run("dequeue empty queue after draining", func(t *testing.T) {
		queue := NewRingQueue[int](2)
		_ = queue.Enqueue(1)
		_, _ = queue.Dequeue()

		_, err := queue.Dequeue()

		if !errors.Is(err, ErrEmptyQueue) {
			t.Errorf("Expected ErrEmptyQueue, got %v", err)
		}
	})
}


