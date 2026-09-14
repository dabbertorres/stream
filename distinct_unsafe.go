//go:build !noseqdistinctunsafe

package stream

import (
	"hash/maphash"
	"unsafe"
)

// DistinctUnsafe returns a Seq that yields only the first occurrence of each
// distinct element of s, in order.
//
// Unlike the free function [Distinct], this method does not require T to be
// comparable: it identifies duplicates by hashing the raw in-memory bytes of
// each element instead of using ==. This makes it unsound for any T whose
// equality isn't fully determined by its bit pattern -- for example, two
// equal strings or slices with different backing arrays will hash
// differently and be treated as distinct, and types containing padding
// bytes, interfaces, maps, or funcs cannot be meaningfully compared this
// way at all. Prefer the free function [Distinct] (which requires T to be
// comparable) whenever T permits it -- the "Unsafe" in this method's name
// is not a formality.
//
// This method can be excluded from the build with the "noseqdistinctunsafe"
// build tag, for callers that don't want the unsafe package linked in.
func (s Seq[T]) DistinctUnsafe() Seq[T] {
	return func(yield func(T) bool) {
		var (
			seen = make(map[uint64]struct{})
			seed = maphash.MakeSeed()
		)

		for v := range s {
			// TODO: it'd be nice to have a better way (read: not using unsafe) to do this
			id := maphash.Bytes(seed,
				unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v)))

			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}

			if !yield(v) {
				return
			}
		}
	}
}
