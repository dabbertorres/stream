package stream

import (
	"iter"
	"slices"
)

func (s Seq[T]) Sorted(less func(lhs, rhs T) bool) Seq[T] {
	sorted := slices.SortedFunc(iter.Seq[T](s), func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})

	return Of(slices.Values(sorted))
}
