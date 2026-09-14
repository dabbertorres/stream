package stream

func (s Seq[T]) Associate[K, V any](by func(T) (K, V)) Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for v := range s {
			k, v := by(v)

			if !yield(Pair[K, V]{k, v}) {
				return
			}
		}
	}
}

// ToMap collects a Seq of [Pair]s into a map, keyed by each Pair's Key.
// If multiple Pairs share the same Key, the last one wins.
func ToMap[K comparable, V any](s Seq[Pair[K, V]]) map[K]V {
	m := make(map[K]V)
	for p := range s {
		m[p.Key] = p.Val
	}
	return m
}
