package stream

func (s Seq[T]) Limit(n int) Seq[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}

		var i int
		for v := range s {
			if !yield(v) {
				return
			}

			i++
			if i >= n {
				return
			}
		}
	}
}
