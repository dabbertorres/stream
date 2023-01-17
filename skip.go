package stream

type skipStream[T any] struct {
	parent streamer[T]
	skip   int
}

func (s skipStream[T]) forEach(f func(T) bool) {
	var total int
	s.parent.forEach(func(value T) bool {
		if total < s.skip {
			total++
			return true
		}
		return f(value)
	})
}

func (s skipStream[T]) capacityHint() int {
	return s.parent.capacityHint() - s.skip
}
