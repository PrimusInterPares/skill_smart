package main

// Порядковый номер задания на курсе: 4
// Номер задачи из задания: 4
// Краткое название: функция, определяющая сбалансированность скобок `(` и `)` в строке
// Time complexity: O(N)
// Memory complexity: O(N)
func IsStringBalanced(s string) bool {
	st := &StackSlice[rune]{}

	for _, char := range s {
		switch char {
		case '(':
			st.Push(char)
		case ')':
			if st.Size() == 0 {
				return false
			}

			_, _ = st.Pop()
		}
	}

	return st.Size() == 0
}

// Номер задачи из задания: 5
// Краткое название: функция, определяющая сбалансированность скобок `(`, `)`, `{`, `}`, `[` и `]` в строке
// Time complexity: O(N)
// Memory complexity: O(N)
func IsStringBalanced2(s string) bool {
	st := &StackSlice[rune]{}

	for _, char := range s {
		switch char {
		case '(', '[', '{':
			st.Push(char)
		case ')':
			if st.Size() == 0 {
				return false
			}

			head, _ := st.Peek()
			if head != '(' {
				return false
			}

			_, _ = st.Pop()
		case ']':
			if st.Size() == 0 {
				return false
			}

			head, _ := st.Peek()
			if head != '[' {
				return false
			}

			_, _ = st.Pop()
		case '}':
			if st.Size() == 0 {
				return false
			}

			head, _ := st.Peek()
			if head != '{' {
				return false
			}

			_, _ = st.Pop()
		}
	}

	return st.Size() == 0
}

// Номер задачи из задания: 6
// Краткое название: функция, возвращающая минимальное значение в стеке за O(1)
// Time complexity: O(1)
// Memory complexity: O(N)
func (st *Stack[T]) Min() (T, error) {
	var zero T
	if len(st.mins) == 0 {
		return zero, ErrEmptyStack
	}

	return st.mins[len(st.mins)-1], nil
}

func (st *Stack[T]) SetLess(less func(a, b T) bool) {
	st.less = less
}

// Номер задачи из задания: 7
// Краткое название: функция, возвращающая среднее значение всех элементов в стеке за O(1)
// Time complexity: O(1)
// Memory complexity: O(1)
func (st *Stack[T]) Avg() (float64, error) {
	if st.Size() == 0 {
		return 0, ErrEmptyStack
	}

	return st.sum / float64(st.Size()), nil
}

// Номер задачи из задания: 8
// Краткое название: функция, которая с помощью двух стеков реализует вычисление постфиксных выражений
// Time complexity: O(N)
// Memory complexity: O(N)
func PostfixCalc(input *Stack[rune]) int {
	numbers := &Stack[int]{}

	for input.Size() > 0 {
		item, _ := input.Pop()

		switch item {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			numbers.Push(int(item - '0'))
		case '+':
			a, _ := numbers.Pop()
			b, _ := numbers.Pop()

			numbers.Push(b + a)
		case '*':
			a, _ := numbers.Pop()
			b, _ := numbers.Pop()

			numbers.Push(b * a)
		case '=':
			res, _ := numbers.Peek()

			return res
		}
	}

	return 0
}

// Рефлексия
// Номер задачи из задания: 9
// Да, я тоже пришёл к тому, что надо просто `стрелки` местами поменять. Элегантное решение.

// Номер задачи из задания: 10
// Видимо, я себе чуть усложнил задачу. Я подумал, что циклы могут быть как вперёд, так и назад,
// поэтому проверял еще и prev ссылки. А вот проверку саму для циклов вперёд сделал через map.
// Просто дойти до конца и проверить, является ли последний элемент хвостом - гениально, как просто! Аж стыдно теперь.

// Номер задачи из задания: 11
// Я подумал, что у пузырька была бы квадратичная сложность даже в среднем случае, и я решил, что такое делать негоже.
// Но и быстро написать какой-нибудь MergeSort я б не смог, поэтому впервые использовал строенную стортировку для slice.
// Она подстраивается под структуру данных и в худшем случае выдаёт O(N*logN)

// Номер задачи из задания: 12
// Я использовал именно описанный вариант: шёл параллельно по обоим спискам,
// собирая 'нужные' значения последовательно в результрующий список.
// Но моё решение, кажись, не сработает для произвольного кол-ва списков,
// т.к. у меня цикл заполнения останавливается, когда один из списко закончился.
// И результирующий список затем дополняется элементами из невыбранного списка.

// Номер задачи из задания: 13
// Классная идея. Надо как-то научиться включать мышление вне рамок, полезно...


