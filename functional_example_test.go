package stream_test

import (
	"fmt"

	"github.com/dabbertorres/stream"
)

type Foo struct {
	X int
}

func (f Foo) Add(y int) int { return f.X + y }

func ExampleChain() {
	var (
		isPositive = func(x int) bool { return x > 0 }
		isEven     = func(x int) bool { return x%2 == 0 }
		doubleIt   = func(x int) int { return x * 2 }
	)

	bizLogic := stream.Chain(
		stream.Chain(
			stream.ApplyRight(stream.Stream[int].Filter, isPositive),
			stream.ApplyRight(stream.Stream[int].Filter, isEven),
		),
		stream.ApplyRight(stream.Stream[int].Transform, doubleIt),
	)

	bizLogic(stream.FromSlice([]int{-3, 4, 1, 12})).
		ForEach(func(i int) { fmt.Println(i) })

	// Output:
	// 8
	// 24
}

// ExamplePipe is the same logic as [ExampleChain], but shows the readability improvement
// of [Pipe] in certain cases.
func ExamplePipe() {
	var (
		isPositive = func(x int) bool { return x > 0 }
		isEven     = func(x int) bool { return x%2 == 0 }
		doubleIt   = func(x int) int { return x * 2 }
	)

	stream.Pipe(
		stream.FromSlice([]int{-3, 4, 1, 12}),
		stream.Chain(
			stream.ApplyRight(stream.Stream[int].Filter, isPositive),
			stream.ApplyRight(stream.Stream[int].Filter, isEven),
		),
		stream.ApplyRight(stream.Stream[int].Transform, doubleIt),
	).ForEach(func(i int) { fmt.Println(i) })

	// Output:
	// 8
	// 24
}

func ExampleApplyLeft() {
	concat := func(prefix, suffix string) string {
		return prefix + " " + suffix
	}

	commonPrefix := stream.ApplyLeft(concat, "INFO:")

	fmt.Println(commonPrefix("it worked!"))

	// Output:
	// INFO: it worked!
}

func ExampleApplyRight() {
	f := stream.ApplyRight(Foo.Add, 5)
	result := f(Foo{X: 3})
	fmt.Println(result)

	// Output:
	// 8
}
