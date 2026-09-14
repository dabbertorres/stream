package stream

func (s Seq[T]) DropWhile(f func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		dropping := true

		for v := range s {
			if dropping {
				if f(v) {
					continue
				}
				dropping = false
			}

			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq[T]) TakeWhile(f func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !f(v) {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}
