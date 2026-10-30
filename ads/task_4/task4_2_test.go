package main

import (
	"errors"
	"math"
	"testing"
)

func Test_IsStringBalanced(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Simple valid", "()", true},
		{"Nested valid", "((()))", true},
		{"Complex valid", "(()((())()))", true},

		{"Single closing", ")", false},
		{"Single opening", "(", false},
		{"Unclosed open brackets", "((()", false},
		{"Extra closing at the end", "())", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStringBalanced(tt.input)

			if result != tt.expected {
				t.Errorf("`%q` balance is %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func Test_IsStringBalanced2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Simple valid", "()[]{}", true},
		{"Nested valid", "{[()]}", true},
		{"Complex valid", "[()]{}{[()()]()}", true},
		
		{"Single closing", ")", false},
		{"Single opening", "[", false},
		{"Wrong closing type", "(]", false},
		{"Interleaved brackets", "([)]", false},
		{"Unclosed open brackets", "({[]", false},
		{"Extra closing at the end", "()[]{}}", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStringBalanced2(tt.input)
			
			if result != tt.expected { t.Errorf("`%q` balance is %v, expected %v", tt.input, result, tt.expected) }
		})
	}
}

func TestStack_Min(t *testing.T) {
	t.Run("In an empty stack", func(t *testing.T) {
		stack := &Stack[int]{}

		stack.SetLess(func(a, b int) bool { return a < b })
		_, err := stack.Min()

		if err == nil { t.Fatal("Error is expected") }
		if !errors.Is(err, ErrEmptyStack) { t.Errorf("ErrEmptyStack is expected, got %v", err) }
	})

	t.Run("Min is a head of a stack", func(t *testing.T) {
		stack := &Stack[int]{}

		stack.SetLess(func(a, b int) bool { return a < b })
		for i := range 10 { stack.Push(i) }

		minVal, err := stack.Min()

		if err != nil { t.Fatalf("unepxected error: %v", err) }
		if minVal != 0 { t.Errorf("Expected 0 as minimum, got %d", minVal) }
	})
	
	t.Run("Min is a tail of a stack", func(t *testing.T) {
		stack := &Stack[int]{}

		stack.SetLess(func(a, b int) bool { return a < b })
		for i := 10; i > 0; i-- { stack.Push(i) }

		minVal, err := stack.Min()

		if err != nil { t.Fatalf("unepxected error: %v", err) }
		if minVal != 1 { t.Errorf("Expected 1 as minimum, got %d", minVal) }
	})

	t.Run("Min is in the middle of a stack", func(t *testing.T) {
		stack := &Stack[int]{}

		stack.SetLess(func(a, b int) bool { return a < b })
		stack.Push(1)
		stack.Push(-10)
		stack.Push(56)

		minVal, err := stack.Min()

		if err != nil { t.Fatalf("unepxected error: %v", err) }
		if minVal != -10 { t.Errorf("Expected -10 as minimum, got %d", minVal) }
	})
}

func floatEquals(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestStack_Avg(t *testing.T) {
	t.Run("Empty stack", func(t *testing.T) {
		st := &Stack[int]{}
	
		_, err := st.Avg()
		if !errors.Is(err, ErrEmptyStack) { t.Errorf("Expected ErrEmptyStack, got %v", err) }
	})

	t.Run("Single element", func(t *testing.T) {
		st := &Stack[int]{}
	
		st.Push(42)

		avg, err := st.Avg()
		if err != nil { t.Fatalf("Unexpected error: %v", err) }
		if !floatEquals(avg, 42.0) { t.Errorf("Expected 42.0, got %f", avg) }
	})

	t.Run("Multiple ints", func(t *testing.T) {
		st := &Stack[int]{}
		
		st.Push(10)
		st.Push(20)
		st.Push(30)

		avg, err := st.Avg()
		if err != nil { t.Fatalf("Unexpected error: %v", err) }
		//nolint:gocritic // (10 + 20 + 30) / 3 = 20
		if !floatEquals(avg, 20.0) { t.Errorf("Expected 20.0, got %f", avg) }
	})

	t.Run("Negative numbers", func(t *testing.T) {
		st := &Stack[int]{}
	
		st.Push(-10)
		st.Push(20)
		st.Push(-30)

		avg, err := st.Avg()
		if err != nil { t.Fatalf("Unexpected error: %v", err) }
		// (-10 + 20 - 30) / 3 = -20 / 3 = -6.666...
		if !floatEquals(avg, -20.0/3.0) { t.Errorf("Expected %f, got %f", -20.0/3.0, avg) }
	})

	t.Run("Push and Pop updates sum", func(t *testing.T) {
		st := &Stack[int]{}
	
		st.Push(10)
		st.Push(20)
		st.Push(30) // Stack: 30, 20, 10. Sum = 60. Avg = 20.0
	
		avg, _ := st.Avg()
		if !floatEquals(avg, 20.0) { t.Errorf("Step 1: Expected 20.0, got %f", avg) }

		st.Pop() // Removes 30. Stack: 20, 10. Sum = 30. Avg = 15.0
		avg, _ = st.Avg()
		if !floatEquals(avg, 15.0) { t.Errorf("Step 2: Expected 15.0, got %f", avg) }

		st.Push(40) // Stack: 40, 20, 10. Sum = 70. Avg = 23.333...
		avg, _ = st.Avg()
		if !floatEquals(avg, 70.0/3.0) { t.Errorf("Step 3: Expected %f, got %f", 70.0/3.0, avg) }
	})

	t.Run("With floats", func(t *testing.T) {
		st := &Stack[float64]{}
		
		st.Push(1.5)
		st.Push(2.5)

		avg, err := st.Avg()
		if err != nil { t.Fatalf("Unexpected error: %v", err) }
		//nolint:gocritic // (1.5 + 2.5) / 2 = 2.0
		if !floatEquals(avg, 2.0) { t.Errorf("Expected 2.0, got %f", avg) }
	})
}

func Test_PostfixCalc(t *testing.T) {
	t.Run("Input stack does not contain `=`", func(t *testing.T) {
		input := &Stack[rune]{}
		
		input.Push('+')
		input.Push('9')
		input.Push('*')
		input.Push('5')
		input.Push('+')
		input.Push('2')
		input.Push('8')

		result := PostfixCalc(input)

		if result != 0 { t.Errorf("Expression must not be calculated, 0 must be returned but got %d", result) }
	})
	
	t.Run("Cmplete expression", func(t *testing.T) {
		input := &Stack[rune]{}
		
		input.Push('=')
		input.Push('+')
		input.Push('9')
		input.Push('*')
		input.Push('5')
		input.Push('+')
		input.Push('2')
		input.Push('8')

		result := PostfixCalc(input)

		if result != 59 { t.Errorf("Expression must eq 59, got %d", result) }
	})

	t.Run("Empty stack", func(t *testing.T) {
		input := &Stack[rune]{}
	
		result := PostfixCalc(input)
		
		if result != 0 { t.Errorf("Expected 0 for empty stack, got %d", result) }
	})

	t.Run("Only equals sign", func(t *testing.T) {
		input := &Stack[rune]{}

		input.Push('=')
		
		result := PostfixCalc(input)
		
		if result != 0 { t.Errorf("Expected 0 for stack with only '=', got %d", result) }
	})

	t.Run("Single number", func(t *testing.T) {
		input := &Stack[rune]{}

		input.Push('=')
		input.Push('7')
		
		result := PostfixCalc(input)
		
		if result != 7 { t.Errorf("Expected 7, got %d", result) }
	})

	t.Run("Missing operands", func(t *testing.T) {
		input := &Stack[rune]{}

		input.Push('=')
		input.Push('+')
		input.Push('8')

		result := PostfixCalc(input)
	
		// Текущая реализация игнорирует ошибки Pop, поэтому вернёт 8 (0 + 8)
		if result != 8 { t.Errorf("Expected 8 (due to ignored errors), got %d", result) }
	})
}


