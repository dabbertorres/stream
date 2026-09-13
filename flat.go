package stream

func (s Seq[T]) FlatMap[O any](mapper func(T) Seq[O]) Seq[O] {
	return func(yield func(O) bool) {
		for v := range s {
			for o := range mapper(v) {
				if !yield(o) {
					return
				}
			}
		}
	}
}
