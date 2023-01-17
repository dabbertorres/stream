package stream

type transformStream[T any] struct {
	parent    streamer[T]
	transform func(T) T
}

func (s transformStream[T]) forEach(f func(T) bool) {
	s.parent.forEach(func(value T) bool {
		return f(s.transform(value))
	})
}

func (f transformStream[T]) capacityHint() int { return f.parent.capacityHint() }
