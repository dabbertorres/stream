package stream

type dropWhileStream[T any] struct {
	parent streamer[T]
	filter func(T) bool
}

func (s dropWhileStream[T]) forEach(f func(T) bool) {
	// dropping...
	s.parent.forEach(func(elem T) bool {
		if s.filter(elem) {
			return true
		}

		f(elem)
		return false
	})

	// pass on everything else
	s.parent.forEach(f)
}

func (s dropWhileStream[T]) capacityHint() int { return s.parent.capacityHint() }

type takeWhileStream[T any] struct {
	parent streamer[T]
	filter func(T) bool
}

func (s takeWhileStream[T]) forEach(f func(T) bool) {
	s.parent.forEach(func(elem T) bool {
		if !s.filter(elem) {
			return false
		}

		return f(elem)
	})

	// drop everything else
}

func (s takeWhileStream[T]) capacityHint() int { return s.parent.capacityHint() }
