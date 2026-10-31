package main

import (
	"errors"
	"testing"
)

func TestQueue_Size(t *testing.T) {
	queue := &Queue[int]{}

	t.Run("Size of an empty queue", func(t *testing.T) {
		size := queue.Size()

		if size != 0 {
			t.Errorf("Size must be 0, got %d", size)
		}
	})

	t.Run("Size of a queue with elements", func(t *testing.T) {
		queue.Enqueue(1)
		queue.Enqueue(2)

		size := queue.Size()

		if size != 2 {
			t.Errorf("Size must be 2, got %d", size)
		}
	})
}

func TestQueue_Dequeue(t *testing.T) {
	queue := &Queue[string]{}

	t.Run("Dequeue from an empty queue", func(t *testing.T) {
		item, err := queue.Dequeue()

		if err == nil {
			t.Error("Error is expected")
		}
		if !errors.Is(err, ErrEmptyQueue) {
			t.Errorf("ErrEmptyQueue is expected, got %v", err)
		}
		var zeroString string
		if item != zeroString {
			t.Errorf("No value must be returned, got %s", item)
		}
	})

	t.Run("Dequeue from a queue with elements", func(t *testing.T) {
		queue.Enqueue("first")
		queue.Enqueue("second")
		queue.Enqueue("third")

		item, err := queue.Dequeue()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item != "first" {
			t.Errorf("Expectes dequeued item to be `first`, got `%s`", item)
		}
	})
}

func TestQueue_Enqueue(t *testing.T) {
	queue := &Queue[string]{}

	t.Run("Equeueing in an empty queue", func(t *testing.T) {
		queue.Enqueue("first")

		if queue.Size() != 1 {
			t.Errorf("Size after enqueueing must become 1, got %d", queue.Size())
		}
	})

	t.Run("Enqueueing in a queue with elements", func(t *testing.T) {
		queue.Enqueue("first")

		oldSize := queue.Size()
		queue.Enqueue("second")

		if queue.Size() != oldSize+1 {
			t.Errorf("Enquque must increase size by 1, increased by %d", oldSize-queue.Size())
		}

		item, err := queue.Dequeue()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if item != "first" {
			t.Errorf("Enqueue must place items into the beginning of the queue")
		}
	})
}
