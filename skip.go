package stream

func (s Seq[T]) Skip(n int) Seq[T] {
	return func(yield func(T) bool) {
		var i int
		for v := range s {
			if i < n {
				i++
				continue
			}

			if !yield(v) {
				return
			}
		}
	}
}
