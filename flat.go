package stream

func Flatten[T any](in Stream[Stream[T]]) Stream[T] {
	return Stream[T]{
		src: flattenStream[T]{
			parent: in.src,
		},
	}
}

type flattenStream[T any] struct {
	parent streamer[Stream[T]]
}

func (s flattenStream[T]) forEach(f func(T) bool) {
	s.parent.forEach(func(stream Stream[T]) bool {
		stream.src.forEach(func(value T) bool {
			return f(value)
		})
		return true
	})
}

func (f flattenStream[T]) capacityHint() int { return f.parent.capacityHint() }
