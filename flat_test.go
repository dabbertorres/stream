package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		actual := Flatten(
			FromSlice([]Stream[int]{
				FromSlice([]int{13, -5, 12}),
				FromSlice([]int{38, 0}),
				FromSlice([]int{121, 53, -300}),
			}),
		).Collect()

		expect := []int{13, -5, 12, 38, 0, 121, 53, -300}
		assert.Equal(t, expect, actual)
	})

	t.Run("child modifiers are respected", func(t *testing.T) {
		actual := Flatten(
			FromSlice([]Stream[int]{
				FromSlice([]int{13, -5, 12}).Filter(func(i int) bool { return i > 0 }),
				FromSlice([]int{38, 0}).Filter(func(i int) bool { return i != 0 }),
				FromSlice([]int{121, 53, -300}).Filter(func(i int) bool { return i < 0 }),
			}),
		).Collect()

		expect := []int{13, 12, 38, -300}
		assert.Equal(t, expect, actual)
	})
}
