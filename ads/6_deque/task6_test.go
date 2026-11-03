package main

import (
	"errors"
	"testing"
)

func TestDeque_AddFront(t *testing.T) {
	t.Run("Add into empty deque", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddFront(3)

		if d.front == nil {
			t.Error("Front block is not allocated")
		}
		if d.front != d.back {
			t.Error("First block must be both front and back")
		}
		if d.front.next != nil {
			t.Error("There must be no next block for a single block")
		}
		if d.front.prev != nil {
			t.Error("There must be no prev block for a single block")
		}
		if d.Size() != 1 {
			t.Errorf("Deque size must be 1, got %d", d.Size())
		}
		v := d.front.items[blockSize-1]
		if v != 3 {
			t.Errorf("Element value must be 63, got %d", v)
		}
	})

	t.Run("Front has free space", func(t *testing.T) {
		d := &Deque[int]{}
		for i := range blockSize - 1 { // from 0 to 62 (63 elements)
			d.AddFront(i)
		}

		d.AddFront(99)

		if d.front != d.back {
			t.Error("Frond must be back (new block is not allocated)")
		}
		if d.front.head != 0 {
			t.Errorf("Index of the front.head must be 0 (full block), got %d", d.front.head)
		}
		if d.Size() != blockSize {
			t.Errorf("Deque size must be %d (full), got %d", blockSize, d.Size())
		}
	})

	t.Run("Front does not have free space (new front block allocation)", func(t *testing.T) {
		d := &Deque[int]{}
		for i := range blockSize { // from 0 to 63 (64 elements) -> full block
			d.AddFront(i)
		}
		oldFront := d.front

		d.AddFront(99) // new block allocation

		if d.front == oldFront {
			t.Error("New front block must have been allocated")
		}
		if d.front.prev != nil {
			t.Error("New front prev must be nil")
		}
		if d.front.next != oldFront {
			t.Error("New front next must be old front")
		}
		if d.back != oldFront {
			t.Error("Old front must become back")
		}
		if d.back.prev != d.front {
			t.Error("Old front prev must be new front")
		}
		if d.back.next != nil {
			t.Error("Old front next must be nil")
		}
	})

	t.Run("Overflow (consecutive pushes)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddFront(i)
		}

		if d.Size() != 200 {
			t.Errorf("Size must be 200, got %d", d.Size())
		}
	})
}

func TestDeque_AddTail(t *testing.T) {
	t.Run("Add into empty deque", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddTail(3)

		if d.back == nil {
			t.Error("Back block is not allocated")
		}
		if d.front != d.back {
			t.Error("First block must be both front and back")
		}
		if d.Size() != 1 {
			t.Errorf("Deque size must be 1, got %d", d.Size())
		}
		v := d.back.items[0]
		if v != 3 {
			t.Errorf("Element value must be 3, got %d", v)
		}
	})

	t.Run("Back has free space", func(t *testing.T) {
		d := &Deque[int]{}
		for i := range blockSize - 1 { // from 0 to 62 (63 elements)
			d.AddTail(i)
		}

		d.AddTail(99)

		if d.front != d.back {
			t.Error("Front must be back (new block is not allocated)")
		}
		if d.back.tail != blockSize {
			t.Errorf("Index of the back.tail must be %d (full block), got %d", blockSize, d.back.tail)
		}
		if d.Size() != blockSize {
			t.Errorf("Deque size must be %d (full), got %d", blockSize, d.Size())
		}
	})

	t.Run("Back does not have free space (new back block allocation)", func(t *testing.T) {
		d := &Deque[int]{}
		for i := range blockSize { // from 0 to 63 (64 elements) -> full block
			d.AddTail(i)
		}
		oldBack := d.back

		d.AddTail(99) // new block allocation

		if d.back == oldBack {
			t.Error("New back block must have been allocated")
		}
		if d.back.next != nil {
			t.Error("New back next must be nil")
		}
		if d.back.prev != oldBack {
			t.Error("New back prev must be old back")
		}
		if d.front != oldBack {
			t.Error("Front must remain old back")
		}
		if d.front.next != d.back {
			t.Error("Front next must be new back")
		}
		if d.front.prev != nil {
			t.Error("Front prev must be nil")
		}
	})

	t.Run("Overflow (consecutive pushes)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddTail(i)
		}

		if d.Size() != 200 {
			t.Errorf("Size must be 200, got %d", d.Size())
		}
	})
}

func TestDeque_RemoveFront(t *testing.T) {
	t.Run("Remove from empty deque", func(t *testing.T) {
		d := &Deque[int]{}

		v, err := d.RemoveFront()

		if !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque is expeceted")
		}
		if v != 0 {
			t.Errorf("Zero value must be returned on error, got %d", v)
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must stay 0, got %d", d.Size())
		}
		if d.front != nil || d.back != nil {
			t.Error("Front and back must stay nil for an untouched deque")
		}
	})

	t.Run("Remove single element", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddFront(3)

		v, err := d.RemoveFront()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 3 {
			t.Errorf("Element value must be 3, got %d", v)
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
		if _, err := d.RemoveFront(); !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque must be returned after the last element is removed")
		}
	})

	t.Run("Front block has remaining elements", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddFront(1)
		d.AddFront(2)
		d.AddFront(3) // items: [_, 3, 2, 1], head = 61, tail = 64
		headBefore := d.front.head

		v, err := d.RemoveFront()

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 3 {
			t.Errorf("Element value must be 3, got %d", v)
		}
		if d.front.head != headBefore+1 {
			t.Errorf("Front head must be %d, got %d", headBefore+1, d.front.head)
		}
		if d.front != d.back {
			t.Error("Front must be back (block must not be unlinked)")
		}
		if d.Size() != 2 {
			t.Errorf("Deque size must be 2, got %d", d.Size())
		}
	})

	t.Run("Front block becomes empty (block must be unlinked)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize { // full first block
			d.AddTail(i)
		}
		d.AddTail(99) // new back block
		oldFront := d.front
		oldBack := d.back

		for i := range blockSize {
			v, err := d.RemoveFront()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}

		if d.front == oldFront {
			t.Error("Empty front block must be unlinked")
		}
		if d.front != oldBack {
			t.Error("Front must become the old back block")
		}
		if d.front.prev != nil {
			t.Error("New front prev must be nil")
		}
		if d.back != oldBack {
			t.Error("Back must remain the old back block")
		}
		if d.Size() != 1 {
			t.Errorf("Deque size must be 1, got %d", d.Size())
		}
		v, err := d.RemoveFront()
		if err != nil || v != 99 {
			t.Errorf("Must return (99, nil), got (%d, %v)", v, err)
		}
	})

	t.Run("New front block becomes empty after overflow", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize { // full block, head = 0
			d.AddFront(i)
		}
		oldFront := d.front

		d.AddFront(99) // new front block with a single element

		v, err := d.RemoveFront()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 99 {
			t.Errorf("Element value must be 99, got %d", v)
		}
		if d.front != oldFront {
			t.Error("Empty front block must be unlinked, old front must be restored")
		}
		if d.front.prev != nil {
			t.Error("Restored front prev must be nil")
		}
		if d.back != oldFront {
			t.Error("Back must be the old front block")
		}
		if d.Size() != blockSize {
			t.Errorf("Deque size must be %d, got %d", blockSize, d.Size())
		}
	})

	t.Run("AddTail + RemoveFront (FIFO)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddTail(i)
		}

		for i := range 200 {
			v, err := d.RemoveFront()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}

		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
	})

	t.Run("Order: AddFront + RemoveFront (LIFO)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddFront(i)
		}

		for i := 199; i >= 0; i-- {
			v, err := d.RemoveFront()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}
	})

	t.Run("Last block is kept when deque becomes empty", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize {
			d.AddTail(i)
		}

		for range 64 {
			if _, err := d.RemoveFront(); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		}
		if d.front == nil || d.back == nil {
			t.Error("Last block must be kept for reuse")
		}
		if d.front != d.back {
			t.Error("Front must be back")
		}
		if d.front.prev != nil || d.front.next != nil {
			t.Error("Single block must have no links")
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
	})

	t.Run("Reuse after full drain", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 100 {
			d.AddTail(i)
		}

		for range 100 {
			d.RemoveFront()
		}

		d.AddFront(1)
		d.AddTail(2)

		if d.Size() != 2 {
			t.Errorf("Deque size must be 2, got %d", d.Size())
		}
		v, err := d.RemoveFront()
		if err != nil || v != 1 {
			t.Errorf("Must return (1, nil), got (%d, %v)", v, err)
		}
		v, err = d.RemoveFront()
		if err != nil || v != 2 {
			t.Errorf("Must return (2, nil), got (%d, %v)", v, err)
		}
	})
}

func TestDeque_RemoveTail(t *testing.T) {
	t.Run("Empty deque", func(t *testing.T) {
		d := &Deque[int]{}

		v, err := d.RemoveTail()

		if !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque expeceted")
		}
		if v != 0 {
			t.Errorf("Zero value must be returned on error, got %d", v)
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must stay 0, got %d", d.Size())
		}
		if d.front != nil || d.back != nil {
			t.Error("Front and back must stay nil for an untouched deque")
		}
	})

	t.Run("Remove single element", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddTail(3)

		v, err := d.RemoveTail()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 3 {
			t.Errorf("Element value must be 3, got %d", v)
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
		if _, err := d.RemoveTail(); !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque must be returned after the last element is removed")
		}
	})

	t.Run("Back block has remaining elements", func(t *testing.T) {
		d := &Deque[int]{}

		d.AddTail(1)
		d.AddTail(2)
		d.AddTail(3) // [1, 2, 3, _], head = 0, tail = 3

		v, err := d.RemoveTail()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 3 {
			t.Errorf("Element value must be 3, got %d", v)
		}
		if d.back.tail != 2 {
			t.Errorf("Index of the back.tail must be 2, got %d", d.back.tail)
		}
		if d.front != d.back {
			t.Error("Front must be back (block must not be unlinked)")
		}
		if d.Size() != 2 {
			t.Errorf("Deque size must be 2, got %d", d.Size())
		}
	})

	t.Run("Back block becomes empty (block must be unlinked)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize { // full first block
			d.AddFront(i)
		}
		d.AddFront(99) // new front block
		oldFront := d.front
		oldBack := d.back

		for i := range blockSize {
			v, err := d.RemoveTail()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}

		if d.back == oldBack {
			t.Error("Empty back block must be unlinked")
		}
		if d.back != oldFront {
			t.Error("Back must become the old front block")
		}
		if d.back.next != nil {
			t.Error("New back next must be nil")
		}
		if d.front != oldFront {
			t.Error("Front must remain the old front block")
		}
		if d.Size() != 1 {
			t.Errorf("Deque size must be 1, got %d", d.Size())
		}
		v, err := d.RemoveTail()
		if err != nil || v != 99 {
			t.Errorf("Must return (99, nil), got (%d, %v)", v, err)
		}
	})

	t.Run("New back block becomes empty after AddTail overflow", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize { // full block, tail = 64
			d.AddTail(i)
		}
		oldBack := d.back
		d.AddTail(99) // new back block with a single element

		v, err := d.RemoveTail()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if v != 99 {
			t.Errorf("Element value must be 99, got %d", v)
		}
		if d.back != oldBack {
			t.Error("Empty back block must be unlinked, old back must be restored")
		}
		if d.back.next != nil {
			t.Error("Restored back next must be nil")
		}
		if d.front != oldBack {
			t.Error("Front must be the old back block")
		}
		if d.Size() != blockSize {
			t.Errorf("Deque size must be %d, got %d", blockSize, d.Size())
		}
	})

	t.Run("Order: AddFront + RemoveTail (FIFO)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddFront(i)
		}

		for i := range 200 {
			v, err := d.RemoveTail()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}

		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
	})

	t.Run("AddTail + RemoveTail (LIFO)", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 200 {
			d.AddTail(i)
		}

		for i := 199; i >= 0; i-- {
			v, err := d.RemoveTail()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
		}
	})

	t.Run("Last block is kept when deque becomes empty", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range blockSize {
			d.AddFront(i)
		}

		for range blockSize {
			if _, err := d.RemoveTail(); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		}

		if d.front == nil || d.back == nil {
			t.Error("Last block must be kept for reuse")
		}
		if d.front != d.back {
			t.Error("Front must be back")
		}
		if d.front.prev != nil || d.front.next != nil {
			t.Error("Single block must have no links")
		}
		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
	})

	t.Run("Reuse after full drain", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 100 {
			d.AddFront(i)
		}

		for range 100 {
			d.RemoveTail()
		}
		d.AddTail(1)
		d.AddFront(2)

		if d.Size() != 2 {
			t.Errorf("Deque size must be 2, got %d", d.Size())
		}
		v, err := d.RemoveTail()
		if err != nil || v != 1 {
			t.Errorf("Must return (1, nil), got (%d, %v)", v, err)
		}
		v, err = d.RemoveTail()
		if err != nil || v != 2 {
			t.Errorf("Must return (2, nil), got (%d, %v)", v, err)
		}
	})
}

func TestDeque_MixedOperations(t *testing.T) {
	t.Run("AddFront + RemoveTail keep a single block alive", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 1000 {
			d.AddFront(i)

			v, err := d.RemoveTail()
			if err != nil {
				t.Fatalf("Unexpected error at %d: %v", i, err)
			}
			if v != i {
				t.Errorf("Element value must be %d, got %d", i, v)
			}
			if d.Size() != 0 {
				t.Errorf("Deque size must be 0, got %d", d.Size())
			}
		}
	})

	t.Run("Interleaved pushes and pops", func(t *testing.T) {
		d := &Deque[int]{}

		for i := range 500 {
			d.AddFront(i)
			d.AddTail(i)
		}

		if d.Size() != 1000 {
			t.Errorf("Deque size must be 1000, got %d", d.Size())
		}

		for range 500 {
			if _, err := d.RemoveFront(); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if _, err := d.RemoveTail(); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		}

		if d.Size() != 0 {
			t.Errorf("Deque size must be 0, got %d", d.Size())
		}
		if _, err := d.RemoveFront(); !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque must be returned for an empty deque")
		}
		if _, err := d.RemoveTail(); !errors.Is(err, ErrEmptyDeque) {
			t.Error("ErrEmptyDeque must be returned for an empty deque")
		}
	})
}
