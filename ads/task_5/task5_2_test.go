package main

import (
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
		{"обычный сдвиг", []int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{"сдвиг на 0", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"n больше размера", []int{1, 2, 3}, 5, []int{3, 1, 2}}, // 5 % 3 = 2
		{"отрицательный n", []int{1, 2, 3, 4}, -13, []int{4, 1, 2, 3}},
		{"очередь пуста", []int{}, 2, []int{}},
		{"один элемент", []int{1}, 10, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue[int]{items: tt.initial}

			q.Rotate(tt.n)

			if !slices.Equal(q.items, tt.expected) {
				t.Errorf("Expected %v, got %v", q.items, tt.expected)
			}
		})
	}
}

func TestQueue_Reverse(t *testing.T) {
	tests := []struct {
		name     string
		initial  []int
		expected []int
	}{
		{"обычная очередь", []int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}},
		{"очередь пуста", []int{}, []int{}},
		{"один элемент", []int{1}, []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queue[int]{items: tt.initial}

			q.Reverse()

			if !slices.Equal(q.items, tt.expected) {
				t.Errorf("Expected %v, got %v", q.items, tt.expected)
			}
		})
	}
}
