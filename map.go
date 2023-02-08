package stream

func Map[In, Out any](in Stream[In], mapper func(In) Out) Stream[Out] {
	return Stream[Out]{
		src: mapStream[In, Out]{
			src:    in.src,
			mapper: mapper,
		},
	}
}

type mapStream[In, Out any] struct {
	src    streamer[In]
	mapper func(In) Out
}

func (m mapStream[In, Out]) forEach(f func(Out) bool) {
	m.src.forEach(func(value In) bool {
		return f(m.mapper(value))
	})
}

func (m mapStream[In, Out]) capacityHint() int { return m.src.capacityHint() }

func FlatMap[In, Out any](in Stream[Stream[In]], mapper func(In) Out) Stream[Out] {
	return Stream[Out]{
		src: flattenStream[Out]{
			parent: mapStream[Stream[In], Stream[Out]]{
				src:    in.src,
				mapper: func(i Stream[In]) Stream[Out] { return Map(i, mapper) },
			},
		},
	}
}
