package stream

import (
	"strconv"
	"testing"
)

// toMapWithSizeHint is an alternative ToMap that first collects to a slice,
// using its length as a size hint for the destination map, to compare
// against ToMap's current make(map[K]V) with no hint.
func toMapWithSizeHint[K comparable, V any](s Seq[Pair[K, V]]) map[K]V {
	elems := s.ToSlice()

	m := make(map[K]V, len(elems))
	for _, p := range elems {
		m[p.Key] = p.Val
	}
	return m
}

func benchmarkPairs(n int) []Pair[int, int] {
	pairs := make([]Pair[int, int], n)
	for i := range pairs {
		pairs[i] = Pair[int, int]{Key: i, Val: i * 2}
	}
	return pairs
}

func BenchmarkToMap(b *testing.B) {
	for _, n := range []int{10, 1_000, 100_000} {
		pairs := benchmarkPairs(n)

		b.Run(strconv.Itoa(n)+"/no_hint", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = ToMap(OfSlice(pairs))
			}
		})

		b.Run(strconv.Itoa(n)+"/size_hint", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = toMapWithSizeHint(OfSlice(pairs))
			}
		})
	}
}
