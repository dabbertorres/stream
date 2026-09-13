package stream

import (
	"iter"
	"slices"
)

func (s Seq[T]) ToSlice() []T {
	return slices.Collect(iter.Seq[T](s))
}

func (s Seq[T]) Append(out []T) []T {
	return slices.AppendSeq(out, iter.Seq[T](s))
}

func (s Seq[T]) All(f func(T) bool) bool {
	for v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

func (s Seq[T]) Any(f func(T) bool) bool {
	for v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func (s Seq[T]) None(f func(T) bool) bool {
	return !s.Any(f)
}

// ToChan drains s into a channel, on a spawned goroutine, so it can be
// used with select or other channel-based code. The channel is closed once
// s is exhausted.
func (s Seq[T]) ToChan() <-chan T {
	ch := make(chan T)

	go func() {
		defer close(ch)
		for v := range s {
			ch <- v
		}
	}()

	return ch
}

func (s Seq[T]) Max(less func(lhs, rhs T) bool) Optional[T] {
	var (
		max  T
		some bool
	)

	for v := range s {
		if !some || less(max, v) {
			max = v
			some = true
		}
	}

	if !some {
		return None[T]()
	}
	return Some(max)
}

func (s Seq[T]) Min(less func(lhs, rhs T) bool) Optional[T] {
	var (
		min  T
		some bool
	)

	for v := range s {
		if !some || less(v, min) {
			min = v
			some = true
		}
	}

	if !some {
		return None[T]()
	}
	return Some(min)
}

func (s Seq[T]) MinMax(less func(lhs, rhs T) bool) (min, max Optional[T]) {
	var (
		minVal, maxVal T
		some           bool
	)

	for v := range s {
		if !some {
			minVal = v
			maxVal = v
			some = true
			continue
		}

		switch {
		case less(v, minVal):
			minVal = v
		case less(maxVal, v):
			maxVal = v
		}
	}

	if !some {
		return None[T](), None[T]()
	}
	return Some(minVal), Some(maxVal)
}
