package stream

func (s Seq[T]) Reduce[O any](start O, reducer func(cumulative O, next T) O) O {
	cumulative := start
	for v := range s {
		cumulative = reducer(cumulative, v)
	}
	return cumulative
}
