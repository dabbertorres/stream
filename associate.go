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

// KeyValueAssociater is a helper to reduce noise when using [Associate] with [KeyValue]s.
type KeyValueAssociater[K comparable, V any] struct {
	src streamer[KeyValue[K, V]]
}

// By creates a map[K]V from the Stream, with keys and values defined by f.
// If f is nil, the identity function is used (K and V without any changes).
func (a KeyValueAssociater[K, V]) By(f func(K, V) (K, V)) map[K]V {
	if f == nil {
		f = func(k K, v V) (K, V) { return k, v }
	}

	m := make(map[K]V, a.src.capacityHint())
	a.src.forEach(func(kv KeyValue[K, V]) bool {
		k, v := f(kv.Key, kv.Val)
		m[k] = v
		return true
	})
	return m
}

// ByTo fills to from the Stream, with keys and values defined by f.
// If f is nil, the identity function is used (K and V without any changes).
func (a KeyValueAssociater[K, V]) ByTo(to map[K]V, f func(K, V) (K, V)) {
	a.src.forEach(func(kv KeyValue[K, V]) bool {
		k, v := f(kv.Key, kv.Val)
		to[k] = v
		return true
	})
}

func AssociateFromKeyValue[K comparable, V any](in Stream[KeyValue[K, V]]) KeyValueAssociater[K, V] {
	return KeyValueAssociater[K, V]{src: in.src}
}
