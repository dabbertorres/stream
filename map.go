package stream

type KeyValue[K comparable, V any] struct {
	Key K
	Val V
}

// KeyValueLess provides a generic [LessFunc] for [KeyValue]s with an [Ordered] key.
func KeyValueLess[K Ordered, V any](lhs, rhs KeyValue[K, V]) bool { return lhs.Key < rhs.Key }

// KeyValueAssociate provides a default function for [KeyValue]s for [Associate].
func KeyValueAssociate[K comparable, V any](kv KeyValue[K, V]) (K, V) { return kv.Key, kv.Val }

func FromMap[K comparable, V any, M ~map[K]V](src M) Stream[KeyValue[K, V]] {
	return Stream[KeyValue[K, V]]{
		src: mapStream[K, V, M]{
			src: src,
		},
	}
}

type mapStream[K comparable, V any, M ~map[K]V] struct{ src M }

func (m mapStream[K, V, M]) forEach(f func(KeyValue[K, V]) bool) {
	for k, v := range m.src {
		if !f(KeyValue[K, V]{Key: k, Val: v}) {
			return
		}
	}
}

func (m mapStream[K, V, M]) capacityHint() int { return len(m.src) }
