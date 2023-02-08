package stream

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
