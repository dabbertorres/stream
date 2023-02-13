package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	ch := make(chan int)
	go func() {
		defer close(ch)
		ch <- 7
		ch <- 11
		ch <- 19
	}()

	var called bool
	actual := Join(
		FromChan(ch),
		FromSlice([]int{42, 23, 9}),
		FromFunc(func() Optional[int] {
			if !called {
				called = true
				return Some(3)
			}
			return None[int]()
		}),
	).Collect()

	expect := []int{7, 11, 19, 42, 23, 9, 3}
	assert.Equal(t, expect, actual)
}
