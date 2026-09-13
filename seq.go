package stream

import "iter"

type Pair[K, V any] struct {
	Key K
	Val V
}

type Seq[T any] iter.Seq[T]

func Of[T any](seq iter.Seq[T]) Seq[T] {
	return Seq[T](seq)
}

func OfChan[C ~(<-chan T), T any](c C) Seq[T] {
	return func(yield func(T) bool) {
		for v := range c {
			if !yield(v) {
				return
			}
		}
	}
}

func OfKeyed[K, V any](seq iter.Seq2[K, V]) Seq[Pair[K, V]] {
	return func(yield func(Pair[K, V]) bool) {
		for k, v := range seq {
			if !yield(Pair[K, V]{k, v}) {
				return
			}
		}
	}
}
