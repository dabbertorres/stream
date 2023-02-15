package stream

type Associater[In any, K comparable, V any] struct {
	src streamer[In]
}

func (a Associater[In, K, V]) By(f func(In) (K, V)) map[K]V {
	m := make(map[K]V, a.src.capacityHint())
	a.ByTo(m, f)
	return m
}

func (a Associater[In, K, V]) ByTo(to map[K]V, f func(In) (K, V)) {
	a.src.forEach(func(in In) bool {
		k, v := f(in)
		to[k] = v
		return true
	})
}

func Associate[K comparable, V any, In any](in Stream[In]) Associater[In, K, V] {
	return Associater[In, K, V]{src: in.src}
}
