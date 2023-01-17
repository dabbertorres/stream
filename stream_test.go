package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Run("zero_initialized", func(t *testing.T) {
		assert.NotPanics(t, func() {
			stream := Stream[int]{}.
				Filter(func(i int) bool { return i%2 == 0 }).
				Distinct().
				Limit(10).
				Skip(2).
				Sorted(func(lhs, rhs int) bool { return true }).
				Transform(func(i int) int { return i * 2 }).
				DropWhile(func(i int) bool { return i < 5 }).
				TakeWhile(func(i int) bool { return i > 5 })

			_ = stream.All(func(i int) bool { return true })
			_ = stream.Any(func(i int) bool { return true })
			_ = stream.None(func(i int) bool { return true })
			_ = stream.Append([]int{})
			_ = stream.Collect()
			_, _ = stream.First()
			_, _ = stream.FirstWhere(func(i int) bool { return true })
			stream.ForEach(func(i int) {})
			_, _ = stream.Max(func(lhs, rhs int) bool { return true })
			_, _ = stream.Min(func(lhs, rhs int) bool { return true })
			_, _, _ = stream.MinMax(func(lhs, rhs int) bool { return true })
			_ = stream.Range()
			stream.RangeTo(make(chan int))
			_ = stream.Reduce(1, func(cumulative, next int) int { return cumulative + next })
		})
	})

	t.Run("call_order", func(t *testing.T) {
		var (
			callOrder int

			filterIndex int
			filterCalls int

			sortedIndex int
			sortedCalls int

			dropWhileIndex int
			dropWhileCalls int

			transformIndex int
			transformCalls int
		)

		actual := FromSlice([]int{-11, 3, 44, 23, 19, -101, 3}).
			Filter(func(i int) bool {
				if filterIndex == 0 {
					callOrder++
					filterIndex = callOrder
				}
				filterCalls++

				return i > 0
			}).
			Distinct().
			Sorted(func(lhs, rhs int) bool {
				if sortedIndex == 0 {
					callOrder++
					sortedIndex = callOrder
				}
				sortedCalls++
				return lhs < rhs
			}).
			Transform(func(i int) int {
				if transformIndex == 0 {
					callOrder++
					transformIndex = callOrder
				}
				transformCalls++
				return i * 2
			}).
			DropWhile(func(i int) bool {
				if dropWhileIndex == 0 {
					callOrder++
					dropWhileIndex = callOrder
				}
				dropWhileCalls++
				return i < 10
			}).
			Limit(1).
			Collect()

		assert.Equal(t, 4, callOrder)

		assert.Equal(t, 1, filterIndex, "Filter() call order")
		assert.Equal(t, 7, filterCalls, "Filter() number of calls")

		assert.Equal(t, 2, sortedIndex, "Sorted() call order")
		assert.NotZero(t, sortedCalls, "Sorted() number of calls") // NOTE: don't check the exact number of calls since it is implementation defined

		assert.Equal(t, 3, transformIndex, "Transform() call order")
		assert.Equal(t, 2, transformCalls, "Transform() number of calls")

		assert.Equal(t, 4, dropWhileIndex, "DropWhile() call order")
		assert.Equal(t, 2, dropWhileCalls, "DropWhile() number of calls")

		expect := []int{38}
		assert.Equal(t, expect, actual)
	})
}

func TestStreamDistinct(t *testing.T) {
	t.Run("ordered", func(t *testing.T) {
		actual := FromSlice([]int{3, 3, 3, 11, 27, 27, 42}).
			Distinct().
			Collect()

		expect := []int{3, 11, 27, 42}
		assert.Equal(t, expect, actual)
	})

	t.Run("unordered", func(t *testing.T) {
		actual := FromSlice([]int{11, 3, 27, 3, 42, 27}).
			Distinct().
			Collect()

		expect := []int{11, 3, 27, 42}
		assert.Equal(t, expect, actual)
	})
}

func TestStreamFilter(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		Filter(func(i int) bool { return i%2 == 0 }).
		Collect()

	expect := []int{8, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamLimit(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		Limit(2).
		Collect()

	expect := []int{3, 8}
	assert.Equal(t, expect, actual)
}

func TestStreamSkip(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		Skip(2).
		Collect()

	expect := []int{11, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamSorted(t *testing.T) {
	t.Run("is_sorted", func(t *testing.T) {
		actual := FromSlice([]int{3, 8, 11, 42}).
			Sorted(func(lhs, rhs int) bool { return lhs < rhs }).
			Collect()

		expect := []int{3, 8, 11, 42}
		assert.Equal(t, expect, actual)
	})

	t.Run("not_sorted", func(t *testing.T) {
		actual := FromSlice([]int{42, 3, 11, 8}).
			Sorted(func(lhs, rhs int) bool { return lhs < rhs }).
			Collect()

		expect := []int{3, 8, 11, 42}
		assert.Equal(t, expect, actual)
	})
}

func TestStreamTransform(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		Transform(func(i int) int { return i * 2 }).
		Collect()

	expect := []int{6, 16, 22, 84}
	assert.Equal(t, expect, actual)
}

func TestStreamDropWhile(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		DropWhile(func(i int) bool { return i < 10 }).
		Collect()

	expect := []int{11, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamTakeWhile(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		TakeWhile(func(i int) bool { return i < 10 }).
		Collect()

	expect := []int{3, 8}
	assert.Equal(t, expect, actual)
}

func TestStreamAll(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		actual := FromSlice([]int{2, 8, 10, 42}).
			All(func(i int) bool { return i%2 == 0 })

		assert.True(t, actual)
	})

	t.Run("false", func(t *testing.T) {
		actual := FromSlice([]int{2, 8, 10, 42}).
			All(func(i int) bool { return i%2 == 1 })

		assert.False(t, actual)
	})
}

func TestStreamAny(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		actual := FromSlice([]int{1, 3, 7, 42}).
			Any(func(i int) bool { return i%2 == 0 })

		assert.True(t, actual)
	})

	t.Run("false", func(t *testing.T) {
		actual := FromSlice([]int{2, 8, 10, 42}).
			All(func(i int) bool { return i%2 == 1 })

		assert.False(t, actual)
	})
}

func TestStreamNone(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		actual := FromSlice([]int{1, 3, 7, 11}).
			None(func(i int) bool { return i%2 == 0 })

		assert.True(t, actual)
	})

	t.Run("false", func(t *testing.T) {
		actual := FromSlice([]int{2, 8, 10, 27}).
			None(func(i int) bool { return i%2 == 1 })

		assert.False(t, actual)
	})
}

func TestStreamAppend(t *testing.T) {
	original := []int{37, 12}

	actual := FromSlice([]int{3, 8, 11, 42}).
		Append(original)

	expect := []int{37, 12, 3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamCollect(t *testing.T) {
	actual := FromSlice([]int{3, 8, 11, 42}).
		Collect()

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamFirst(t *testing.T) {
	t.Run("stream is not empty", func(t *testing.T) {
		result, found := FromSlice([]int{3, 8, 11, 42}).
			First()
		if assert.True(t, found) {
			assert.Equal(t, 3, result)
		}
	})

	t.Run("stream is empty", func(t *testing.T) {
		result, found := FromSlice[int, []int](nil).
			First()
		if assert.False(t, found) {
			assert.Equal(t, 0, result)
		}
	})
}

func TestStreamFirstIf(t *testing.T) {
	result, found := FromSlice([]int{3, 8, 11, 42}).
		FirstWhere(func(i int) bool { return i > 10 })

	if assert.True(t, found) {
		assert.Equal(t, 11, result)
	}
}

func TestStreamForEach(t *testing.T) {
	var actual []int
	FromSlice([]int{3, 8, 11, 42}).
		ForEach(func(i int) { actual = append(actual, i) })

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestStreamMax(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		actual, found := FromSlice([]int{11, 8, 42, 3}).
			Max(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.True(t, found) {
			assert.Equal(t, 42, actual)
		}
	})

	t.Run("empty", func(t *testing.T) {
		actual, found := FromSlice([]int{}).
			Max(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.False(t, found) {
			assert.Equal(t, 0, actual)
		}
	})
}

func TestStreamMin(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		actual, found := FromSlice([]int{11, 8, 42, 3}).
			Min(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.True(t, found) {
			assert.Equal(t, 3, actual)
		}
	})

	t.Run("empty", func(t *testing.T) {
		actual, found := FromSlice([]int{}).
			Min(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.False(t, found) {
			assert.Equal(t, 0, actual)
		}
	})
}

func TestStreamMinMax(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		actualMin, actualMax, found := FromSlice([]int{11, 8, 42, 3}).
			MinMax(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.True(t, found) {
			assert.Equal(t, 3, actualMin)
			assert.Equal(t, 42, actualMax)
		}
	})

	t.Run("empty", func(t *testing.T) {
		actualMin, actualMax, found := FromSlice([]int{}).
			MinMax(func(lhs, rhs int) bool { return lhs < rhs })

		if assert.False(t, found) {
			assert.Equal(t, 0, actualMin)
			assert.Equal(t, 0, actualMax)
		}
	})
}

func TestStreamRange(t *testing.T) {
	ch := FromSlice([]int{3, 8, 11, 42}).Range()

	assert.Equal(t, 3, <-ch)
	assert.Equal(t, 8, <-ch)
	assert.Equal(t, 11, <-ch)
	assert.Equal(t, 42, <-ch)
	assert.Empty(t, ch)
}

func TestStreamRangeTo(t *testing.T) {
	ch := make(chan int)
	go FromSlice([]int{3, 8, 11, 42}).RangeTo(ch)

	assert.Equal(t, 3, <-ch)
	assert.Equal(t, 8, <-ch)
	assert.Equal(t, 11, <-ch)
	assert.Equal(t, 42, <-ch)
	assert.Empty(t, ch)
}

func TestStreamReduce(t *testing.T) {
	actual := FromSlice([]int{3, 8, 8, 11, 42}).
		Distinct().
		Reduce(1, func(cumulative, next int) int { return cumulative + next })

	expect := 1 + 3 + 8 + 11 + 42
	assert.Equal(t, expect, actual)
}

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

	fmt.Println(filterCalled)    // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(transformCalled) // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(stream.First())  // NOTE: the stream has now been consumed
	fmt.Println(stream.FirstWhere(func(i int) bool { return i%2 == 1 }))
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
	fmt.Println(stream.First())
	fmt.Println(stream.FirstWhere(func(i int) bool { return i%2 == 1 }))
	fmt.Println(stream.First())
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

func ExampleFromFunc() {
	var transformCalled bool

	i := 1
	stream := FromFunc(func() (int, bool) {
		v := i
		i *= 2
		return v, v <= 64
	}).Skip(1).
		Limit(10).
		Transform(func(elem int) int { transformCalled = true; return elem * 3 })

	fmt.Println(transformCalled) // NOTE: false is printed here - streams are lazily evaluated
	fmt.Println(stream.Collect())
	fmt.Println(transformCalled)

	// Output:
	// false
	// [6 12 24 48 96 192]
	// true
}

func ExampleFromDecoder() {
	type bar struct {
		Foo int `json:"foo"`
	}

	data := []byte(`
    {"foo": 7}
    {"foo": 5}
    {"foo": 7}
    {"foo": -13}
`)

	results := FromDecoder[bar](json.NewDecoder(bytes.NewReader(data)), nil).
		Distinct().
		Filter(func(b bar) bool { return b.Foo >= 0 }).
		Collect()

	fmt.Println(results)

	// Output:
	// [{7} {5}]
}
