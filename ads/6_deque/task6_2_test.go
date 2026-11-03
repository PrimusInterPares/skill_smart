package main

import (
	"errors"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	t.Run("Empty and single character cases", func(t *testing.T) {
		if !IsPalindrome("") {
			t.Error("Empty string must be a palindrome")
		}

		if !IsPalindrome("a") {
			t.Error("Single character string 'a' must be a palindrome")
		}

		if !IsPalindrome("я") {
			t.Error("Single Unicode character string 'я' must be a palindrome")
		}
	})

	t.Run("Valid palindromes (ASCII)", func(t *testing.T) {
		tests := []struct {
			name string
			str  string
		}{
			{name: "Even length", str: "abba"},
			{name: "Odd length", str: "racecar"},
			{name: "Repeated characters", str: "vvvvvv"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if !IsPalindrome(tc.str) {
					t.Errorf("String '%s' must be recognized as a palindrome", tc.str)
				}
			})
		}
	})

	t.Run("Valid palindromes (Unicode/Cyrillic)", func(t *testing.T) {
		// Важные тесты для Go, так как обычный обход по индексам s[i]
		// ломается на многобайтовых символах (рунах)
		tests := []struct {
			name string
			str  string
		}{
			{name: "Cyrillic odd length", str: "шалаш"},
			{name: "Cyrillic even length", str: "довод"},
			{name: "Cyrillic complex symmetric", str: "топот"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if !IsPalindrome(tc.str) {
					t.Errorf("Unicode string '%s' must be recognized as a palindrome", tc.str)
				}
			})
		}
	})

	t.Run("Invalid palindromes (Asymmetric)", func(t *testing.T) {
		tests := []struct {
			name string
			str  string
		}{
			{name: "Completely different", str: "golang"},
			{name: "Different only at edges", str: "racecab"},
			{name: "Different only in center", str: "rabacar"},
			{name: "Case sensitive difference", str: "Anna"}, // Если ваш метод чувствителен к регистру
			{name: "Cyrillic asymmetric", str: "палиндром"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				if IsPalindrome(tc.str) {
					t.Errorf("String '%s' is not a palindrome, but got true", tc.str)
				}
			})
		}
	})

	t.Run("Long strings (Stress test)", func(t *testing.T) {
		// Генерируем длинную строку-палиндром (аналог вашего теста на Overflow)
		// 1000 символов 'a' + 'b' + 1000 символов 'a'
		size := 1000
		buf := make([]byte, size*2+1)
		for i := range size {
			buf[i] = 'a'
			buf[size*2-i] = 'a'
		}
		buf[size] = 'b'
		longPalindrome := string(buf)

		if !IsPalindrome(longPalindrome) {
			t.Error("Long generated string (2001 chars) must be a palindrome")
		}

		// Ломаем палиндром в самом конце
		buf[size*2] = 'z'
		longInvalid := string(buf)

		if IsPalindrome(longInvalid) {
			t.Error("Long asymmetric string must not be recognized as a palindrome")
		}
	})
}

func TestDeque_Min(t *testing.T) {
	tests := []struct {
		name        string
		elems       []int
		expectedMin int
		expectedErr error
	}{
		{"Empty deque", []int{}, 0, ErrEmptyDeque},
		{"Single element", []int{42}, 42, nil},
		{"Positive numbers", []int{1, 2, 3}, 1, nil},
		{"Negative numbers", []int{-1, -2, -3}, -3, nil},
		{"Mixed positive and negative numbers", []int{-1, 2, -4, 0, 5}, -4, nil},
		{"Identical numbers", []int{7, 7, 7}, 7, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deque[int]{}

			d.SetLess(func(a, b int) bool { return a < b })
			for _, i := range tt.elems {
				d.AddTail(i)
			}

			actual, err := d.Min()

			if err != tt.expectedErr {
				t.Errorf("Unexpected error: %v", err)
			}
			if actual != tt.expectedMin {
				t.Errorf("Expected min is %v, got %d", tt.expectedMin, actual)
			}
		})
	}

	t.Run("With strings", func(t *testing.T) {
		d := &Deque[string]{}

		d.SetLess(func(a, b string) bool { return a < b })

		d.AddFront("orange")
		d.AddTail("apple")
		d.AddTail("banana")

		actual, err := d.Min()
		expected := "apple"

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if actual != expected {
			t.Errorf("Expected min is `%s`, got `%s`", actual, expected)
		}
	})

	t.Run("Missing Less function", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddTail(5)
		d.AddTail(2)

		_, err := d.Min()
		if !errors.Is(err, ErrNoLessFunc) {
			t.Errorf("Expected ErrNoLessFunc, got %v", err)
		}
	})
}

func TestDynArrayDeque_AddFront(t *testing.T) {
	t.Run("Basic LIFO order", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddFront(1)
		d.AddFront(2)
		d.AddFront(3)

		if d.items.count != 3 {
			t.Errorf("expected count 3, got %d", d.items.count)
		}

		val, _ := d.RemoveHead()
		if val != 3 {
			t.Errorf("expected 3, got %d", val)
		}
	})
}

func TestDynArrayDeque_RemoveFront(t *testing.T) {
	t.Run("Empty deque error", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		_, err := d.RemoveHead()
		if !errors.Is(err, ErrEmptyDeque) {
			t.Errorf("expected error on empty deque")
		}
	})

	t.Run("FIFO order after AddTail", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddTail(1)
		d.AddTail(2)
		d.AddTail(3)

		val, err := d.RemoveHead()
		if err != nil || val != 1 {
			t.Errorf("expected 1, nil, got %d, %v", val, err)
		}

		val, err = d.RemoveHead()
		if err != nil || val != 2 {
			t.Errorf("expected 2, nil, got %d, %v", val, err)
		}
	})
}

func TestDynArrayDeque_AddTail(t *testing.T) {
	t.Run("Basic FIFO order", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddTail(1)
		d.AddTail(2)
		d.AddTail(3)

		if d.items.count != 3 {
			t.Errorf("expected count 3, got %d", d.items.count)
		}

		val, _ := d.RemoveHead()
		if val != 1 {
			t.Errorf("expected 1, got %d", val)
		}
	})
}

func TestDynArrayDeque_RemoveTail(t *testing.T) {
	t.Run("Empty deque", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		_, err := d.RemoveTail()
		if !errors.Is(err, ErrEmptyDeque) {
			t.Errorf("expected ErrEmptyDeque on empty deque")
		}
	})

	t.Run("Sequential RemoveTail", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddFront(1)
		d.AddFront(2)
		d.AddFront(3) // [3, 2, 1]

		val, err := d.RemoveTail()
		if err != nil || val != 1 {
			t.Errorf("expected 1, nil, got %d, %v", val, err)
		}
		val, err = d.RemoveTail()
		if err != nil || val != 2 {
			t.Errorf("expected 2, nil, got %d, %v", val, err)
		}
	})

	t.Run("Remove the only element and add again", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddTail(100)
		d.RemoveTail()
		d.AddTail(99)

		val, err := d.RemoveTail()
		if err != nil || val != 99 {
			t.Errorf("expected 99, nil, got %d, %v", val, err)
		}
	})
}

func TestDynArrayDeque_WrapAroundStress(t *testing.T) {
	t.Run("Head and tail cross physical boundaries without resize", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init() // Capacity is 16

		// initial head (8) -> (8 - 10 + 16) % 16 = 14
		for i := 1; i <= 10; i++ {
			d.AddFront(i) // [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
		}

		for range 5 {
			_, err := d.RemoveHead() // head moves to (14 + 5) % 16 = 3, [5, 4, 3, 2, 1]
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		for i := 11; i <= 18; i++ {
			d.AddTail(i) // [5, 4, 3, 2, 1, 11, 12, 13, 14, 15, 16, 17, 18]
		}

		if d.items.count != 13 {
			t.Errorf("expected count 13, got %d", d.items.count)
		}

		val, _ := d.RemoveHead()
		if val != 5 {
			t.Errorf("expected 5, got %d", val)
		}

		val, _ = d.RemoveTail()
		if val != 18 {
			t.Errorf("expected 18, got %d", val)
		}
	})
}

func TestDynArrayDeque_ShrinkToMinimumBoundary(t *testing.T) {
	t.Run("Capacity never drops below minimalBufferSize", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		for i := range 20 { // increase buffer to 32
			d.AddTail(i)
		}
		if d.items.capacity != 32 {
			t.Errorf("expected capacity 32 after 20 adds, got %d", d.items.capacity)
		}

		for range 15 { // buffer decrease on 15 to 21 and on 10 to 16
			d.RemoveHead()
		}

		if d.items.capacity != 16 {
			t.Errorf("expected capacity to hit minimalBufferSize 16, got %d", d.items.capacity)
		}

		d.AddTail(99)
		val, err := d.RemoveTail()
		if err != nil || val != 99 {
			t.Errorf("expected 99, nil, got %d, %v", val, err)
		}
	})
}

func TestDynArrayDeque_AlternatingOperations(t *testing.T) {
	t.Run("Mixed sequence maintains integrity", func(t *testing.T) {
		d := &DynArrayDeque[int]{items: &DynArray[int]{}}
		d.Init()

		d.AddTail(1)  // [1]
		d.AddFront(2) // [2, 1]
		d.AddTail(3)  // [2, 1, 3]
		d.AddFront(4) // [4, 2, 1, 3]

		val, _ := d.RemoveHead() // [2, 1, 3]
		if val != 4 {
			t.Errorf("expected 4, got %d", val)
		}

		d.AddTail(5)            // [2, 1, 3, 5]
		val, _ = d.RemoveTail() // [2, 1, 3]
		if val != 5 {
			t.Errorf("expected 5, got %d", val)
		}

		val, _ = d.RemoveHead() // [1, 3]
		if val != 2 {
			t.Errorf("expected 2, got %d", val)
		}

		val, _ = d.RemoveTail() // [1]
		if val != 3 {
			t.Errorf("expected 3, got %d", val)
		}

		if d.items.count != 1 {
			t.Errorf("expected count 1, got %d", d.items.count)
		}

		val, _ = d.RemoveHead() // []
		if val != 1 {
			t.Errorf("expected 1, got %d", val)
		}

		if d.items.count != 0 {
			t.Errorf("expected count 0, got %d", d.items.count)
		}
	})
}


