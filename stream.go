package stream

type streamer[T any] interface {
	forEach(f func(T) bool)
	capacityHint() int
}

type Stream[T any] struct {
	src streamer[T]
}

func (s Stream[T]) Distinct() Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: distinctStream[T]{
			parent: s.src,
		},
	}
}

func (s Stream[T]) Filter(filter func(T) bool) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: filterStream[T]{
			parent: s.src,
			filter: filter,
		},
	}
}

func (s Stream[T]) Limit(n int) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: limitStream[T]{
			parent: s.src,
			limit:  n,
		},
	}
}

func (s Stream[T]) Skip(n int) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: skipStream[T]{
			parent: s.src,
			skip:   n,
		},
	}
}

func (s Stream[T]) Sorted(less LessFunc[T]) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: sortedStream[T]{
			parent: s.src,
			less:   less,
		},
	}
}

func (s Stream[T]) Transform(f func(T) T) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: transformStream[T]{
			parent:    s.src,
			transform: f,
		},
	}
}

func (s Stream[T]) DropWhile(f func(T) bool) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: dropWhileStream[T]{
			parent: s.src,
			filter: f,
		},
	}
}

func (s Stream[T]) TakeWhile(f func(T) bool) Stream[T] {
	if s.src == nil {
		return s
	}

	return Stream[T]{
		src: takeWhileStream[T]{
			parent: s.src,
			filter: f,
		},
	}
}

func (s Stream[T]) All(f func(T) bool) (allMatched bool) {
	if s.src == nil {
		return false
	}

	allMatched = true
	s.src.forEach(func(elem T) bool {
		if !f(elem) {
			allMatched = false
			return false
		}
		return true
	})

	return allMatched
}

func (s Stream[T]) Any(f func(T) bool) (anyMatched bool) {
	if s.src == nil {
		return false
	}

	s.src.forEach(func(elem T) bool {
		if f(elem) {
			anyMatched = true
			return false
		}
		return true
	})

	return anyMatched
}

func (s Stream[T]) None(f func(T) bool) (noneMatched bool) {
	if s.src == nil {
		return true // none matched if empty...
	}

	noneMatched = true
	s.src.forEach(func(elem T) bool {
		if f(elem) {
			noneMatched = false
			return false
		}
		return true
	})

	return noneMatched
}

func (s Stream[T]) Append(out []T) []T {
	if s.src == nil {
		return out
	}

	s.src.forEach(func(value T) bool {
		out = append(out, value)
		return true
	})
	return out
}

func (s Stream[T]) Collect() (out []T) {
	if s.src == nil {
		return nil
	}

	if hint := s.src.capacityHint(); hint > 0 {
		out = make([]T, 0, hint)
	}

	s.src.forEach(func(value T) bool {
		out = append(out, value)
		return true
	})
	return out
}

func (s Stream[T]) First() (result T, ok bool) {
	if s.src == nil {
		return result, false
	}

	s.src.forEach(func(value T) bool {
		result = value
		ok = true
		return false
	})
	return result, ok
}

func (s Stream[T]) FirstWhere(f func(T) bool) (result T, ok bool) {
	if s.src == nil {
		return result, false
	}

	s.src.forEach(func(value T) bool {
		if f(value) {
			result = value
			ok = true
			return false
		}
		return true
	})
	return result, ok
}

func (s Stream[T]) ForEach(f func(T)) {
	if s.src == nil {
		return
	}

	s.src.forEach(func(value T) bool {
		f(value)
		return true
	})
}

func (s Stream[T]) Max(less LessFunc[T]) (max T, found bool) {
	if s.src == nil {
		return max, false
	}

	// set the first value
	s.src.forEach(func(elem T) bool {
		max = elem
		found = true
		return false
	})

	s.src.forEach(func(elem T) bool {
		if less(max, elem) {
			max = elem
		}
		return true
	})

	return max, found
}

func (s Stream[T]) Min(less LessFunc[T]) (min T, found bool) {
	if s.src == nil {
		return min, false
	}

	// set the first value
	s.src.forEach(func(elem T) bool {
		min = elem
		found = true
		return false
	})

	s.src.forEach(func(elem T) bool {
		if less(elem, min) {
			min = elem
		}
		return true
	})

	return min, found
}

func (s Stream[T]) MinMax(less LessFunc[T]) (min, max T, found bool) {
	if s.src == nil {
		return min, max, false
	}

	// set the first value
	s.src.forEach(func(elem T) bool {
		max = elem
		min = elem
		found = true
		return false
	})

	s.src.forEach(func(elem T) bool {
		switch {
		case less(elem, min):
			min = elem
		case less(max, elem):
			max = elem
		}

		return true
	})

	return min, max, found
}

func (s Stream[T]) Range() <-chan T {
	ch := make(chan T)

	if s.src == nil {
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)
		s.src.forEach(func(value T) bool {
			ch <- value
			return true
		})
	}()

	return ch
}

func (s Stream[T]) RangeTo(ch chan<- T) {
	defer close(ch)

	if s.src == nil {
		return
	}

	s.src.forEach(func(value T) bool {
		ch <- value
		return true
	})
}

func (s Stream[T]) Reduce(start T, reducer func(cumulative T, next T) T) T {
	if s.src == nil {
		return start
	}

	cumulative := start
	s.src.forEach(func(value T) bool {
		cumulative = reducer(cumulative, value)
		return true
	})
	return cumulative
}
