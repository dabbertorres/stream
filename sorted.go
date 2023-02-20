package stream

import "container/heap"

func Sorted[T any](in Stream[T], less LessFunc[T]) Stream[T] {
	return Stream[T]{
		src: sortedStream[T]{
			parent: in.src,
			less:   less,
		},
	}
}

func SortedBy[T any](less LessFunc[T]) func(Stream[T]) Stream[T] {
	return ApplyRight(Sorted[T], less)
}

type sortedStream[T any] struct {
	parent streamer[T]
	less   LessFunc[T]
}

func (s sortedStream[T]) forEach(f func(T) bool) {
	data := genericHeap[T]{less: s.less}
	if hint := s.parent.capacityHint(); hint > 0 {
		data.list = make([]T, 0, hint)
	}

	s.parent.forEach(func(elem T) bool {
		heap.Push(&data, elem)
		return true
	})

	for data.Len() > 0 {
		elem := heap.Pop(&data).(T)
		if !f(elem) {
			return
		}
	}
}

func (s sortedStream[T]) capacityHint() int { return s.parent.capacityHint() }

type genericHeap[T any] struct {
	list []T
	less LessFunc[T]
}

func (h genericHeap[T]) Len() int           { return len(h.list) }
func (h genericHeap[T]) Less(i, j int) bool { return h.less(h.list[i], h.list[j]) }
func (h genericHeap[T]) Swap(i, j int)      { h.list[i], h.list[j] = h.list[j], h.list[i] }
func (h *genericHeap[T]) Push(x any)        { h.list = append(h.list, x.(T)) }
func (h *genericHeap[T]) Pop() any {
	idx := len(h.list) - 1
	end := h.list[idx]
	h.list = h.list[:idx]
	return end
}
