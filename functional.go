package stream

// Chain creates a unary function In -> Out from two functions such that In -> Mid -> Out.
func Chain[Out, Mid, In any, First ~func(In) Mid, Second ~func(Mid) Out](first First, second Second) func(In) Out {
	return func(in In) Out {
		return second(first(in))
	}
}

// Pipe applies in to first, passes its result to second, and returns second's result.
func Pipe[Out, Mid, In any, First ~func(In) Mid, Second ~func(Mid) Out](in In, first First, second Second) Out {
	return Chain(first, second)(in)
}

// ApplyLeft takes a binary function (Left, Right) -> Out, and a value for it's left argument, then returns
// a unary function Right -> Out.
func ApplyLeft[Left, Right, Out any, Func ~func(Left, Right) Out](f Func, left Left) func(Right) Out {
	return func(right Right) Out {
		return f(left, right)
	}
}

// ApplyRight takes a binary function, (Left, Right) -> Out, and a value for it's right argument, then returns
// a unary function Left -> Out.
func ApplyRight[Left, Right, Out any, Func ~func(Left, Right) Out](f Func, right Right) func(Left) Out {
	return func(left Left) Out {
		return f(left, right)
	}
}
