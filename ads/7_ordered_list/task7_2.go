package main

import ()

func (l *OrderedList[T]) RemoveDuplicates() {
	if l.Count() < 2 {
		return
	}

	node := l.head
	next := node.next
	for next != nil {
		if l.Compare(next.value, node.value) == 0 {
			node.next = next.next

			if next.next != nil {
				next.next.prev = node
			} else {
				l.tail = node
			}

			l.count--

			next = node.next
		} else {
			node = node.next
			next = node.next
		}
	}
}
