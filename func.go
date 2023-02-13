package stream

func FromFunc[T any](f func() Optional[T]) Stream[T] {
	return Stream[T]{src: funcStream[T]{f: f}}
}

type funcStream[T any] struct{ f func() Optional[T] }

func (s funcStream[T]) forEach(f func(T) bool) {
	for {
		val, ok := s.f().Get()
		if !ok {
			return
		}

		if !f(val) {
			return
		}
	}
}

func (s funcStream[T]) capacityHint() int { return 0 }
