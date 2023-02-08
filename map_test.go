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
		FromSlice([]int{5, 10, 2, 90}),
		func(x int) Stream[bool] { return FromSlice([]bool{x%2 == 0, x > 5}) },
	).Collect()

	expect := []bool{false, false, true, true, true, false, true, true}
	assert.Equal(t, expect, actual)
}
