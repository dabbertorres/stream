package stream

func Map[In, Out any](in Stream[In], mapper func(In) Out) Stream[Out] {
	return Stream[Out]{
		src: mapperStream[In, Out]{
			src:    in.src,
			mapper: mapper,
		},
	}
}

func MapBy[In, Out any](mapper func(In) Out) func(Stream[In]) Stream[Out] {
	return ApplyRight(Map[In, Out], mapper)
}

func FlatMap[In, Out any](in Stream[In], mapper func(In) Stream[Out]) Stream[Out] {
	return Stream[Out]{
		src: flattenStream[Out]{
			parent: mapperStream[In, Stream[Out]]{
				src:    in.src,
				mapper: mapper,
			},
		},
	}
}

func FlatMapBy[In, Out any](mapper func(In) Stream[Out]) func(Stream[In]) Stream[Out] {
	return ApplyRight(FlatMap[In, Out], mapper)
}

type mapperStream[In, Out any] struct {
	src    streamer[In]
	mapper func(In) Out
}

func (m mapperStream[In, Out]) forEach(f func(Out) bool) {
	m.src.forEach(func(value In) bool {
		return f(m.mapper(value))
	})
}

func (m mapperStream[In, Out]) capacityHint() int { return m.src.capacityHint() }
