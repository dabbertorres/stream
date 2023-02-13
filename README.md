# stream

streams package for Go.

[docs](https://pkg.go.dev/github.com/dabbertorres/stream)

## Quick Examples

### Stream a slice

```go
func ExampleFromSlice() {
	var (
		filterCalled    bool
		transformCalled bool
	)

	stream := FromSlice([]int{3, 8, 11, 24, 37, 42}).
		Filter(func(elem int) bool { filterCalled = true; return elem >= 10 }).
		Skip(1).
		Limit(2).
		Transform(func(elem int) int { transformCalled = true; return elem * 2 })

	fmt.Println(filterCalled)         // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(transformCalled)      // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(stream.First().Get()) // NOTE: the stream has now been consumed
	fmt.Println(stream.FirstWhere(func(i int) bool { return i%2 == 1 }).Get())
	fmt.Println(filterCalled)
	fmt.Println(transformCalled)

	// Output:
	// false
	// false
	// 48 true
	// 0 false
	// true
	// true
}
```

### Stream a channel

```go
func ExampleFromChan() {
	var (
		filterCalled    bool
		transformCalled bool
	)

	ch := make(chan int)
	go func() {
		defer close(ch)
		ch <- 3
		ch <- 8
		ch <- 11
		ch <- 24
		ch <- 37
		ch <- 42
	}()

	stream := FromChan(ch).
		Filter(func(elem int) bool { filterCalled = true; return elem >= 10 }).
		Skip(1).
		Limit(2).
		Transform(func(elem int) int { transformCalled = true; return elem * 2 })

	fmt.Println(filterCalled)    // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(transformCalled) // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(stream.First().Get())
	fmt.Println(stream.FirstWhere(func(i int) bool { return i%2 == 1 }).Get())
	fmt.Println(stream.First().Get())
	fmt.Println(filterCalled)
	fmt.Println(transformCalled)

	// Output:
	// false
	// false
	// 48 true
	// 0 false
	// 0 false
	// true
	// true
}
```
