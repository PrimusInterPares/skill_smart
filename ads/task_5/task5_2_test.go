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
