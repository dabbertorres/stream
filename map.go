package stream

func (s Seq[T]) Map[O any](mapper func(T) O) Seq[O] {
	return func(yield func(O) bool) {
		for v := range s {
			if !yield(mapper(v)) {
				return
			}
		}
	}
}
