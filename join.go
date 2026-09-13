package stream

// Join returns a Seq that yields all of each of seqs' elements, in order.
func Join[T any](seqs ...Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		for _, s := range seqs {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}
