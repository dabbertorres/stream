package stream

func Join[T any](streams ...Stream[T]) Stream[T] {
	join := joinStream[T]{
		parents: make([]streamer[T], len(streams)),
	}

	for i, s := range streams {
		join.parents[i] = s.src
	}

	return Stream[T]{src: join}
}

type joinStream[T any] struct {
	parents []streamer[T]
}

func (s joinStream[T]) forEach(f func(T) bool) {
	for _, parent := range s.parents {
		parent.forEach(f)
	}
}

func (s joinStream[T]) capacityHint() (total int) {
	for _, parent := range s.parents {
		total += parent.capacityHint()
	}
	return total
}
