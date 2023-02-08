package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	actual := Map(
		FromSlice([]int{5, 10, 3, 90}),
		func(x int) bool { return x%2 == 0 },
	).Collect()

	expect := []bool{false, true, false, true}
	assert.Equal(t, expect, actual)
}

func TestFlatMap(t *testing.T) {
	actual := FlatMap(
		FromSlice([]Stream[int]{
			FromSlice([]int{5, 10, 3, 90}),
			FromSlice([]int{37, 2, -5, 15}),
		}),
		func(x int) bool { return x%2 == 0 },
	).Collect()

	expect := []bool{false, true, false, true, false, true, false, false}
	assert.Equal(t, expect, actual)
}
