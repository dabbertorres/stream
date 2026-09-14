package stream

// Distinct returns a Seq that yields only the first occurrence of each
// distinct element of s, in order, using == to compare elements.
//
// See also [Seq.DistinctUnsafe] for a fluent version that does not require T
// to be comparable, at the cost of a weaker (unsafe) notion of equality.
func Distinct[T comparable](s Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{})

		for v := range s {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}

			if !yield(v) {
				return
			}
		}
	}
}
