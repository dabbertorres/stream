package stream

import (
	"bytes"
	"encoding/json"
)

type Optional[T any] struct {
	value T
	some  bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		some:  true,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{some: false}
}

func FromOptional[T any](opt Optional[T]) Stream[T] {
	return Stream[T]{src: opt}
}

func (o Optional[T]) Some() bool { return o.some }
func (o Optional[T]) None() bool { return !o.some }

func (o Optional[T]) Get() (T, bool) { return o.value, o.some }
func (o Optional[T]) MustGet() T {
	if o.some {
		return o.value
	}
	panic("Optional[T] is none")
}

func (o Optional[T]) GetOrDefault(defaultVal T) T {
	if o.some {
		return o.value
	}
	return defaultVal
}

func (o Optional[T]) GetOrDefaultFunc(defaultFunc func() T) T {
	if o.some {
		return o.value
	}
	return defaultFunc()
}

func (o Optional[T]) IfSome(f func(T)) {
	if o.some {
		f(o.value)
	}
}

func (o Optional[T]) IfNone(f func()) {
	if !o.some {
		f()
	}
}

func (o Optional[T]) forEach(f func(T) bool) {
	if o.some {
		_ = f(o.value)
	}
}

func (o Optional[T]) capacityHint() int {
	if o.some {
		return 1
	}
	return 0
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.some {
		return json.Marshal(o.value)
	}

	return []byte("null"), nil
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.some = false
	if bytes.Equal([]byte("null"), data) {
		return nil
	}

	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}

	o.some = true
	return nil
}
