package main

func IsPalindrome(s string) bool {
	if len(s) <= 1 {
		return true
	}

	deque := &Deque[rune]{}
	for _, r := range s {
		deque.AddFront(r)
	}

	for deque.Size() > 1 {
		left, _ := deque.RemoveFront()
		right, _ := deque.RemoveTail()

		if left != right {
			return false
		}
	}

	return true
}
