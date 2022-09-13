package hash_set

type HashSet[T comparable] map[T]struct{}

func (h HashSet[T]) Insert(item T) {
	h[item] = struct{}{}
}

func (h HashSet[T]) Size() int {
	return len(h)
}
