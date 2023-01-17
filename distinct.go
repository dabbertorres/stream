package stream

import (
	"hash/maphash"
	"unsafe"
)

type distinctStream[T any] struct {
	parent streamer[T]
}

func (s distinctStream[T]) forEach(f func(T) bool) {
	var (
		uniqueElems = make(map[uint64]struct{})
		seed        = maphash.MakeSeed()
	)

	s.parent.forEach(func(elem T) bool {
		// TODO: it'd be nice to have a better way (read: not using unsafe) to do this
		id := maphash.Bytes(seed,
			unsafe.Slice((*byte)(unsafe.Pointer(&elem)), unsafe.Sizeof(elem)))

		if _, ok := uniqueElems[id]; !ok {
			// distinct!
			uniqueElems[id] = struct{}{}
			f(elem)
		}

		return true
	})
}

func (s distinctStream[T]) capacityHint() int { return s.parent.capacityHint() }

func Identity[T comparable](v T) T { return v }
