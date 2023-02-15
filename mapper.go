package stream

type Mapper[In, Out any] struct {
	src streamer[In]
}

func (m Mapper[In, Out]) By(mapper func(In) Out) Stream[Out] {
	return Stream[Out]{
		src: mapperStream[In, Out]{
			src:    m.src,
			mapper: mapper,
		},
	}
}

func Map[Out, In any](in Stream[In]) Mapper[In, Out] {
	return Mapper[In, Out]{
		src: in.src,
	}
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

type FlatMapper[In, Out any] struct {
	src streamer[In]
}

func FlatMap[Out, In any](in Stream[In]) FlatMapper[In, Out] {
	return FlatMapper[In, Out]{
		src: in.src,
	}
}

func (m FlatMapper[In, Out]) By(mapper func(In) Stream[Out]) Stream[Out] {
	return Stream[Out]{
		src: flattenStream[Out]{
			parent: mapperStream[In, Stream[Out]]{
				src:    m.src,
				mapper: mapper,
			},
		},
	}
}
