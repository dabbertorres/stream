package stream

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

type Ordered interface {
	Number | ~string
}

// LessFunc is a function that returns true when lhs is less than rhs.
type LessFunc[T any] func(lhs, rhs T) bool

// OrderedLess provides a generic [LessFunc] for builtin types that are [Ordered].
func OrderedLess[T Ordered](lhs, rhs T) bool { return lhs < rhs }

func min[N Number](x, y N) N {
	if x < y {
		return x
	}
	return y
}
