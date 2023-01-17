package stream

import (
	"errors"
	"io"
)

// Decoder is implemented by most of the encoding/* packages in the standard library,
// and often by third-party packages providing similar capabilities.
type Decoder interface {
	Decode(any) error
}

// FromDecoder reads T values from the given Decoder until an error is encountered,
// which results in the end of the Stream.
// If the error is not io.EOF, onError is called with the error.
func FromDecoder[T any](dec Decoder, onError func(error)) Stream[T] {
	if onError == nil {
		onError = func(error) {}
	}

	return Stream[T]{
		src: decoderStream[T]{
			dec:   dec,
			onErr: onError,
		},
	}
}

type decoderStream[T any] struct {
	dec   Decoder
	onErr func(error)
}

func (s decoderStream[T]) forEach(f func(T) bool) {
	for {
		var next T
		if err := s.dec.Decode(&next); err != nil {
			if !errors.Is(err, io.EOF) {
				s.onErr(err)
			}
			return
		}

		if !f(next) {
			return
		}
	}
}

func (s decoderStream[T]) capacityHint() int { return 0 }
