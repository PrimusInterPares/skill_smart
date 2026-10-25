package main

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

// Номер задачи из задания: 8
// Рефлексия
// Т.к. длины списков заранее равные, исключение не делал. В остальном я пришёл к такому же алгоритму.
// Давно не работал со списками, пришлось перенастроить мышление.


