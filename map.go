package stream

type KeyValue[K comparable, V any] struct {
	Key K
	Val V
}

// KeyValueLess provides a generic [LessFunc] for [KeyValue]s with an [Ordered] key.
func KeyValueLess[K Ordered, V any](lhs, rhs KeyValue[K, V]) bool { return lhs.Key < rhs.Key }

// AssociateKeyValue provides a default function for [KeyValue]s for [Associate].
func AssociateKeyValue[K comparable, V any](kv KeyValue[K, V]) (K, V) { return kv.Key, kv.Val }

// ValueTransform is a helper wrapper to simplify transforming only the value of a [KeyValue].
func ValueTransform[K comparable, V any](f func(V) V) func(KeyValue[K, V]) KeyValue[K, V] {
	return func(kv KeyValue[K, V]) KeyValue[K, V] {
		return KeyValue[K, V]{
			Key: kv.Key,
			Val: f(kv.Val),
		}
	}
}

// ValueMapper is a helper wrapper to simplify mapping only the value of a [KeyValue].
func ValueMapper[OutV any, K comparable, InV any](f func(InV) OutV) func(KeyValue[K, InV]) KeyValue[K, OutV] {
	return func(kv KeyValue[K, InV]) KeyValue[K, OutV] {
		return KeyValue[K, OutV]{
			Key: kv.Key,
			Val: f(kv.Val),
		}
	}
}

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
