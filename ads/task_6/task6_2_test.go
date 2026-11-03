package main

import (
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
}
