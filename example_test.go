package stream

import "fmt"

func ExampleOfSlice() {
	var (
		filterCalled bool
		mapCalled    bool
	)

	seq := OfSlice([]int{3, 8, 11, 24, 37, 42}).
		Filter(func(elem int) bool { filterCalled = true; return elem >= 10 }).
		Skip(1).
		Limit(2).
		Map(func(elem int) int { mapCalled = true; return elem * 2 })

	fmt.Println(filterCalled)      // NOTE: false is printed here - Seqs are lazily evaluated
	fmt.Println(mapCalled)         // NOTE: false is printed here - Seqs are lazily evaluated
	fmt.Println(seq.First().Get()) // NOTE: the Seq has now been partially consumed
	fmt.Println(seq.FirstWhere(func(i int) bool { return i%2 == 1 }).Get())
	fmt.Println(filterCalled)
	fmt.Println(mapCalled)

	// Output:
	// false
	// false
	// 48 true
	// 0 false
	// true
	// true
}

func ExampleOfChan() {
	var (
		filterCalled bool
		mapCalled    bool
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

	seq := OfChan[<-chan int](ch).
		Filter(func(elem int) bool { filterCalled = true; return elem >= 10 }).
		Skip(1).
		Limit(2).
		Map(func(elem int) int { mapCalled = true; return elem * 2 })

	fmt.Println(filterCalled) // NOTE: false is printed here - Seqs are lazily evaluated
	fmt.Println(mapCalled)    // NOTE: false is printed here - Seqs are lazily evaluated
	fmt.Println(seq.First().Get())
	fmt.Println(seq.FirstWhere(func(i int) bool { return i%2 == 1 }).Get())
	fmt.Println(seq.First().Get())
	fmt.Println(filterCalled)
	fmt.Println(mapCalled)

	// Output:
	// false
	// false
	// 48 true
	// 0 false
	// 0 false
	// true
	// true
}
