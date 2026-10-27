package main

import(
	"fmt"
)

// Порядковый номер задания на курсе: 3
// Номер задачи из задания: 6
// Краткое название: DynArray на основе банковского метода
const operationPrice int = 3 // усреднённая стоимость операции вставки/удаления - 3

type BankDynArray[T any] struct {
	DynArray[T]
	operationAccount int // уменьшается на стоимость реаллокации, увеличивается на стоимость операции
}

func NewBankDynArray[T any]() *BankDynArray[T] {
	bda := &BankDynArray[T]{
		operationAccount: 0,
	}
	bda.Init()

	return bda
}

func (bda *BankDynArray[T]) MakeArray(sz int) {
	bda.operationAccount -= bda.capacity
  bda.DynArray.MakeArray(sz)

	// проверка для отладки доказательства методом банков
	// if bda.operationAccount < 0 {
	//	panic(fmt.Sprintf("Bankrupt! Account: %d", bda.operationAccount))
	// }
}

func (bda *BankDynArray[T]) Append(itm T) {
	bda.DynArray.Append(itm)
	bda.operationAccount += operationPrice
}

func (bda *BankDynArray[T]) Insert(itm T, index int) error {
	err := bda.DynArray.Insert(itm, index)

	if err == nil {
		bda.operationAccount += operationPrice
	}

	return err
}

func (bda *BankDynArray[T]) Remove(index int) error {
	err := bda.DynArray.Remove(index)

	if err == nil {
		bda.operationAccount += operationPrice
	}

	return err
}

// Номер задачи из задания: 7
// Краткое название: многомерный DynArray
type MultiDynArray[T any] struct {
	depth int
	nodes DynArray[*MultiDynArray[T]] // hierarchy
	values DynArray[T]
}

func NewMultiDynArray[T any](depth int) (*MultiDynArray[T], error) {
	if depth <= 0 {
		var dummy MultiDynArray[T]
		return &dummy, fmt.Errorf("depth must be positive, got: %d", depth)
	}

	return buildMultiDynArray[T](depth), nil
}

func buildMultiDynArray[T any](depth int) *MultiDynArray[T] {
	if depth == 1 {
		leaf := &MultiDynArray[T]{depth: depth}
		leaf.values.Init()

		return leaf
	}

	node := &MultiDynArray[T]{depth: depth}
	node.nodes.Init()
	
	return node
}

// GetItem возвращает элемент. Ожидает ровно a.depth индексов.
func (mda *MultiDynArray[T]) GetItem(indices ...int) (T, error) {
	var dummy T
	if len(indices) != mda.depth {
		return dummy, fmt.Errorf("wrong number of indices: %d, expected: %d", len(indices), mda.depth)
	}

	// 1D массив
	if mda.depth == 1 {
		return mda.values.GetItem(indices[0])
	}

	// спускаемся к дочернему узлу
	child, err := mda.nodes.GetItem(indices[0])
	if err != nil {
		return dummy, err
	}

	return child.GetItem(indices[1:]...)
}

func (mda *MultiDynArray[T]) Append(value T, indices ...int) error {
    if len(indices) != mda.depth - 1 {
        return fmt.Errorf("expected %d indices, got %d", mda.depth - 1, len(indices))
    }

    if mda.depth == 1 {
        mda.values.Append(value)

        return nil
    }

    index := indices[0]

    if index >= mda.nodes.count {
			for mda.nodes.count <= index {
				child := buildMultiDynArray[T](mda.depth - 1)
        mda.nodes.Append(child)
      }
    }

    child, _ := mda.nodes.GetItem(index)
	
    return child.Append(value, indices[1:]...)
}

func (mda *MultiDynArray[T]) Insert(value T, indices ...int) error {
    if len(indices) != mda.depth {
        return fmt.Errorf("wrong number of indices: expected %d, got %d", mda.depth, len(indices))
    }

    if mda.depth == 1 {
        return mda.values.Insert(value, indices[0])
    }

    index := indices[0]
    if index >= mda.nodes.count {
        return fmt.Errorf("index %d is out of bounds: count is %d", index, mda.nodes.count)
    }

    child, _ := mda.nodes.GetItem(index)
    return child.Insert(value, indices[1:]...)
}

// Номер задачи из задания: 8
// Рефлексия
// Т.к. длины списков заранее равные, исключение не делал. В остальном я пришёл к такому же алгоритму.
// Давно не работал со списками, пришлось перенастроить мышление.


