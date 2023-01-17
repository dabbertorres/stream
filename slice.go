package stream

func FromSlice[T any, S ~[]T](src S) Stream[T] {
	return Stream[T]{src: &sliceStream[T, S]{src: src}}
}

type sliceStream[T any, S ~[]T] struct{ src S }

func (s *sliceStream[T, S]) forEach(f func(T) bool) {
	for len(s.src) > 0 {
		v := s.src[0]
		s.src = s.src[1:]
		if !f(v) {
			return
		}
	}
}

func (s sliceStream[T, S]) capacityHint() int { return len(s.src) }
