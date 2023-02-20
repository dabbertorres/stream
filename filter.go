package stream

type filterStream[T any] struct {
	parent streamer[T]
	filter func(T) bool
}

func (s filterStream[T]) forEach(f func(T) bool) {
	s.parent.forEach(func(value T) bool {
		if s.filter(value) {
			return f(value)
		}
		return true
	})
}

func (f filterStream[T]) capacityHint() int { return f.parent.capacityHint() }

func Filter[T any](in Stream[T], filter func(T) bool) Stream[T] {
	return Stream[T]{
		src: filterStream[T]{
			parent: in.src,
			filter: filter,
		},
	}
}

func FilterBy[T any](filter func(T) bool) func(Stream[T]) Stream[T] {
	return ApplyRight(Filter[T], filter)
}
