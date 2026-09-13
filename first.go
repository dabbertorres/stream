package stream

func (s Seq[T]) First() Optional[T] {
	for v := range s {
		return Some(v)
	}
	return None[T]()
}

func (s Seq[T]) FirstWhere(f func(T) bool) Optional[T] {
	for v := range s {
		if f(v) {
			return Some(v)
		}
	}
	return None[T]()
}
