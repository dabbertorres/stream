package stream

import (
	"errors"
	"io"
	"maps"
	"slices"
)

func OfSlice[T any, S ~[]T](src S) Seq[T] {
	return Of(slices.Values(src))
}

func OfMap[K comparable, V any, M ~map[K]V](src M) Seq[Pair[K, V]] {
	return OfKeyed(maps.All(src))
}

func OfOptional[T any](opt Optional[T]) Seq[T] {
	return func(yield func(T) bool) {
		if v, ok := opt.Get(); ok {
			yield(v)
		}
	}
}

func OfFunc[T any](f func() Optional[T]) Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, ok := f().Get()
			if !ok {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}

// Decoder is implemented by most of the encoding/* packages in the standard library,
// and often by third-party packages providing similar capabilities.
type Decoder interface {
	Decode(any) error
}

// OfDecoder reads T values from the given Decoder until an error is encountered,
// which results in the end of the [Seq].
// If the error is not io.EOF, onError is called with the error.
// onError may be nil, in which case the error is silently dropped.
func OfDecoder[T any](dec Decoder, onError func(error)) Seq[T] {
	if onError == nil {
		onError = func(error) {}
	}

	return func(yield func(T) bool) {
		for {
			var next T
			if err := dec.Decode(&next); err != nil {
				if !errors.Is(err, io.EOF) {
					onError(err)
				}
				return
			}

			if !yield(next) {
				return
			}
		}
	}
}
