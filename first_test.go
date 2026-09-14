package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqFirst(t *testing.T) {
	t.Run("non_empty", func(t *testing.T) {
		v, ok := OfSlice([]int{3, 8, 11}).First().Get()
		assert.True(t, ok)
		assert.Equal(t, 3, v)
	})

	t.Run("empty", func(t *testing.T) {
		_, ok := OfSlice([]int{}).First().Get()
		assert.False(t, ok)
	})

	t.Run("does_not_pull_past_first", func(t *testing.T) {
		var pulled []int
		src := Seq[int](func(yield func(int) bool) {
			for _, v := range []int{3, 8, 11} {
				pulled = append(pulled, v)
				if !yield(v) {
					return
				}
			}
		})

		src.First()

		assert.Equal(t, []int{3}, pulled)
	})
}

func TestSeqFirstWhere(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		v, ok := OfSlice([]int{3, 8, 11}).FirstWhere(func(i int) bool { return i%2 == 0 }).Get()
		assert.True(t, ok)
		assert.Equal(t, 8, v)
	})

	t.Run("not_found", func(t *testing.T) {
		_, ok := OfSlice([]int{3, 9, 11}).FirstWhere(func(i int) bool { return i%2 == 0 }).Get()
		assert.False(t, ok)
	})
}
