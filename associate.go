package stream

func Associate[In any, K comparable, V any](
	in Stream[In],
	keyValFunc func(In) (K, V),
) map[K]V {
	out := make(map[K]V, in.src.capacityHint())
	AssociateTo(in, keyValFunc, out)
	return out
}

func AssociateTo[In any, K comparable, V any](
	in Stream[In],
	keyValFunc func(In) (K, V),
	m map[K]V,
) {
	in.src.forEach(func(elem In) bool {
		k, v := keyValFunc(elem)
		m[k] = v
		return true
	})
}
