package main

import(
	"errors"
	"testing"
	"unsafe"
)

func TestStack_Size(t *testing.T) {
	t.Run("When a stack is empty", func(t *testing.T) {
		stack := &Stack[int]{}

		if stack.Size() != 0 { t.Error("Empty stack size must eq 0") }
	})

	t.Run("A stack has elements", func(t *testing.T) {
		stack := &Stack[int]{}
		stack.Push(1)
		stack.Push(2)
		stack.Push(3)

		if stack.Size() != 3 { t.Errorf("Stack size 3 is expected, got %d", stack.Size()) }
	})
}

func TestStack_Peek(t *testing.T) {
	t.Run("Peek from an empty stack", func(t *testing.T) {
		stack := &Stack[int]{}

		_, err := stack.Peek()

		if err == nil { t.Fatal("Expected empty stack error") }
		if !errors.Is(err, ErrEmptyStack) { t.Errorf("Got error %v, expected %v", err, ErrEmptyStack) }
	})

	t.Run("Peek from a stack with elements", func(t *testing.T) {
		stack := &Stack[int]{}

		stack.Push(21)
		stack.Push(111)

		item, err := stack.Peek()
		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if stack.Size() != 2 { t.Errorf("Stack size must be 2, got %d", stack.Size()) }
		if item != 111 { t.Errorf("Head value must be 111, got %d", item) }
	})
}

func TestStack_Pop(t *testing.T) {
	t.Run("Pop from an empty stack", func(t *testing.T) {
		stack := &Stack[int]{}

		_, err := stack.Pop()

		if err == nil { t.Error("Expected error") }
		if !errors.Is(err, ErrEmptyStack) { t.Errorf("Got error %v, expected %v", err, ErrEmptyStack) }
	})

	t.Run("Pop from a stack with elements", func(t *testing.T) {
		stack := &Stack[string]{}

		stack.Push("one")
		stack.Push("two")
		stack.Push("three")
		stack.Push("drink")

		item, err := stack.Pop()

		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if item != "drink" { t.Errorf("Expected `drink` value, got `%s`", item) }

		head, err := stack.Peek()
		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if head != "three" { t.Errorf("Expected `three` as new head, got `%s`", head) }
	})

	t.Run("Pop until empty", func(t *testing.T) {
    stack := &Stack[int]{}

		stack.Push(1)
    stack.Push(2)
    _, _ = stack.Pop()
    _, _ = stack.Pop()

		_, err := stack.Pop()
    if !errors.Is(err, ErrEmptyStack) {
	    t.Error("Expected ErrEmptyStack after full drain")
    }
    if stack.Size() != 0 { t.Errorf("Size must be 0") }
	})

}

func TestStack_Push(t *testing.T) {
	t.Run("Push into an empty stack", func(t *testing.T) {
		stack := &Stack[int]{}

		if stack.Size() != 0 { t.Fatalf("Stack size of 0 is expected, got %d", stack.Size()) }

		stack.Push(1)

		if stack.Size() != 1 { t.Errorf("Stack size after push must be 1, got %d", stack.Size()) }
		
		item, err := stack.Peek()
		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if item != 1 { t.Errorf("Last element of the stack must be 1, got %d", item) }
	})

	t.Run("Push into a stack with elements", func(t *testing.T) {
		stack := &Stack[int]{}
		stack.Push(1)
		stack.Push(2)

		if stack.Size() != 2 { t.Fatalf("Stack size of 2 is expected, got %d", stack.Size()) }

		stack.Push(11)

		if stack.Size() != 3 { t.Errorf("Stack size after push must be 3, got %d", stack.Size()) }

		item, err := stack.Peek()
		if err != nil { t.Fatalf("unexpected error: %v", err) }
		if item != 11 { t.Errorf("Last element of the stack must be 11, got %d", item) }
	})
}

func TestStack_LIFO(t *testing.T) {
	stack := &Stack[int]{}

	for i := 1; i <= 5; i++ { stack.Push(i) }
  for i := 5; i >= 1; i-- {
  	item, _ := stack.Pop()
    if item != i { t.Errorf("Expected %d, got %d", i, item) }
  }
}

func TestStack_AnyTypes(t *testing.T) {
	st := &Stack[any]{}
	
	st.Push(42)
  st.Push("hello")
  st.Push(3.14)
	
	if st.Size() != 3 { t.Errorf("Expected size 3, got %d", st.Size()) }

	item, err := st.Pop()
  if err != nil { t.Fatalf("unexpected error: %v", err) }
	if f, ok := item.(float64); !ok || f != 3.14 { t.Errorf("Expected 3.14, got %v", item) }
}

func TestStackSlice(t *testing.T) {
	t.Run("Zero obj for Garbage Collector", func(t *testing.T) {
    type Big struct { data [1024]byte }
    stack := &StackSlice[*Big]{}

		item := &Big{}
		item.data[0] = 42
    stack.Push(item)
    stack.Pop()

		if cap(stack.items) == 0 { t.Skip("No backing array to check") }

    type sliceHeader struct {
			Data unsafe.Pointer
			Len int
			Cap int
		}

		// Исходный slice: len = 0, cap = 1, backing array = [nil]
		//                                       ^
		//                                      здесь был item, обнулённый плосле Pop
		header:= (*sliceHeader)(unsafe.Pointer(&stack.items))
		newHeader := sliceHeader{
    	Data: header.Data,
      Len:  header.Cap, // 1
      Cap:  header.Cap, // 1
		}

		// *(*TargetType)(unsafe.Pointer(&source)) - Возьми адрес source, интерпретируй его как указатель на TargetType,
		// и разыменуй — получи значение типа TargetType.
		//
		// (*[]*Big)(...) интерпретирует ... как указатель на slice []*Big,
		// *(...) разыменовывает ..., возвращая сам slice []*Big
		newSlice := *(*[]*Big)(unsafe.Pointer(&newHeader))

    // Проверяем, что элемент на позиции 0 был обнулен
    if newSlice[0] != nil { t.Errorf("Popped slot must be nil to allow GC collection, got %v", newSlice[0]) }
	})
	
	t.Run("No shrink when cap <= 1024", func(t *testing.T) {
    stack := &StackSlice[int]{}

    for i := range 500 { stack.Push(i) }

		initialCap := cap(stack.items)
    for range 490 { stack.Pop() } // len < cap/4

		if cap(stack.items) != initialCap { t.Errorf("Small stack must not shrink") }
  })

  t.Run("Shrink when len < cap/4 and cap > 1024", func(t *testing.T) {
    stack := &StackSlice[int]{}

    for i := range 2000 { stack.Push(i) } // cap >= 2000
    for range 1900 { stack.Pop() }   // len=100, cap/4=500

		if cap(stack.items) >= 2000 { t.Errorf("Stack must shrink, got cap=%d", cap(stack.items)) }
    if stack.Size() != 100 { t.Errorf("Data must be preserved after shrink") }
  })

  t.Run("No shrink at exact boundary len == cap/4", func(t *testing.T) {
		stack := &StackSlice[int]{}

		for i := range 2000 { stack.Push(i) }

		// Подгоняем len ровно к cap/4
    for len(stack.items) > cap(stack.items) / 4 + 1 { stack.Pop() }

		capBefore := cap(stack.items)
    stack.Pop() // теперь len == cap/4

		if cap(stack.items) != capBefore { t.Errorf("At boundary len == cap / 4 shrink must NOT trigger") }
  })
}


