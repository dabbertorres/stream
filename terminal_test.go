package stream

import (
	"cmp"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeqToSlice(t *testing.T) {
	actual := OfSlice([]int{3, 8, 11, 42}).ToSlice()

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestSeqAppend(t *testing.T) {
	out := []int{1, 2}
	actual := OfSlice([]int{3, 8, 11}).Append(out)

	expect := []int{1, 2, 3, 8, 11}
	assert.Equal(t, expect, actual)
}

func TestSeqAll(t *testing.T) {
	t.Run("all_match", func(t *testing.T) {
		actual := OfSlice([]int{2, 4, 6}).All(func(i int) bool { return i%2 == 0 })
		assert.True(t, actual)
	})

	t.Run("not_all_match", func(t *testing.T) {
		actual := OfSlice([]int{2, 4, 5}).All(func(i int) bool { return i%2 == 0 })
		assert.False(t, actual)
	})

	t.Run("empty", func(t *testing.T) {
		actual := OfSlice([]int{}).All(func(i int) bool { return false })
		assert.True(t, actual)
	})

	t.Run("short_circuits", func(t *testing.T) {
		var calls int
		OfSlice([]int{2, 5, 4}).All(func(i int) bool {
			calls++
			return i%2 == 0
		})
		assert.Equal(t, 2, calls)
	})
}

func TestSeqAny(t *testing.T) {
	t.Run("some_match", func(t *testing.T) {
		actual := OfSlice([]int{1, 3, 4}).Any(func(i int) bool { return i%2 == 0 })
		assert.True(t, actual)
	})

	t.Run("none_match", func(t *testing.T) {
		actual := OfSlice([]int{1, 3, 5}).Any(func(i int) bool { return i%2 == 0 })
		assert.False(t, actual)
	})

	t.Run("short_circuits", func(t *testing.T) {
		var calls int
		OfSlice([]int{1, 4, 3}).Any(func(i int) bool {
			calls++
			return i%2 == 0
		})
		assert.Equal(t, 2, calls)
	})
}

func TestSeqNone(t *testing.T) {
	t.Run("none_match", func(t *testing.T) {
		actual := OfSlice([]int{1, 3, 5}).None(func(i int) bool { return i%2 == 0 })
		assert.True(t, actual)
	})

	t.Run("some_match", func(t *testing.T) {
		actual := OfSlice([]int{1, 3, 4}).None(func(i int) bool { return i%2 == 0 })
		assert.False(t, actual)
	})
}

func TestSeqToChan(t *testing.T) {
	ch := OfSlice([]int{3, 8, 11}).ToChan(context.Background())

	var actual []int
	for v := range ch {
		actual = append(actual, v)
	}

	expect := []int{3, 8, 11}
	assert.Equal(t, expect, actual)
}

func TestSeqToChanCancel(t *testing.T) {
	var n int
	src := OfFunc(func() Optional[int] {
		n++
		return Some(n)
	})

	ctx, cancel := context.WithCancel(context.Background())
	ch := src.ToChan(ctx)

	<-ch
	cancel()

	select {
	case _, ok := <-ch:
		assert.False(t, ok, "ToChan's channel should be closed after ctx is canceled")
	case <-time.After(time.Second):
		t.Fatal("ToChan's goroutine did not exit after ctx was canceled")
	}
}

func TestSeqMax(t *testing.T) {
	t.Run("non_empty", func(t *testing.T) {
		v, ok := OfSlice([]int{3, 42, 8, 11}).Max(cmp.Less).Get()
		assert.True(t, ok)
		assert.Equal(t, 42, v)
	})

	t.Run("empty", func(t *testing.T) {
		_, ok := OfSlice([]int{}).Max(cmp.Less).Get()
		assert.False(t, ok)
	})
}

func TestSeqMin(t *testing.T) {
	t.Run("non_empty", func(t *testing.T) {
		v, ok := OfSlice([]int{3, 42, -8, 11}).Min(cmp.Less).Get()
		assert.True(t, ok)
		assert.Equal(t, -8, v)
	})

	t.Run("empty", func(t *testing.T) {
		_, ok := OfSlice([]int{}).Min(cmp.Less).Get()
		assert.False(t, ok)
	})
}

func TestSeqMinMax(t *testing.T) {
	t.Run("non_empty", func(t *testing.T) {
		min, max := OfSlice([]int{3, 42, -8, 11}).MinMax(cmp.Less)

		minVal, minOk := min.Get()
		assert.True(t, minOk)
		assert.Equal(t, -8, minVal)

		maxVal, maxOk := max.Get()
		assert.True(t, maxOk)
		assert.Equal(t, 42, maxVal)
	})

	t.Run("empty", func(t *testing.T) {
		min, max := OfSlice([]int{}).MinMax(cmp.Less)

		assert.True(t, min.None())
		assert.True(t, max.None())
	})
}
