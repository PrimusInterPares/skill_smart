package main

import "testing"

// 1. Базовый тест
func TestAdd_Basic(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

// 2. Table-Driven Test (Золотой стандарт в Go)
func TestAdd_Table(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -1, -1, -2},
		{"with zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// 3. Тестирование ошибок (стандартными средствами)
func TestDivide(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err) // Прерывает тест, т.к. дальше проверять бессмысленно
		}
		if got != 5.0 {
			t.Errorf("Divide(10, 2) = %f; want 5.0", got)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Fatal("expected an error, but got nil")
		}
		if err.Error() != "division by zero" {
			t.Errorf("expected 'division by zero', got '%s'", err.Error())
		}
	})
}
