package stream

type limitStream[T any] struct {
	parent streamer[T]
	limit  int
}

func (s limitStream[T]) forEach(f func(T) bool) {
	var total int
	s.parent.forEach(func(value T) bool {
		notDone := f(value)
		total++
		return notDone && total < s.limit
	})
}

func (s limitStream[T]) capacityHint() int {
	return min(s.limit, s.parent.capacityHint())
}
