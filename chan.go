package stream

func FromChan[T any](src <-chan T) Stream[T] {
	return Stream[T]{src: chanStream[T]{src: src}}
}

type chanStream[T any] struct{ src <-chan T }

func (c chanStream[T]) forEach(f func(T) bool) {
	for v := range c.src {
		if !f(v) {
			return
		}
	}
}

func (c chanStream[T]) capacityHint() int { return 0 }
