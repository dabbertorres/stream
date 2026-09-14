package stream

func (s Seq[T]) Filter(filter func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !filter(v) {
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}
