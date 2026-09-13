package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqLimit(t *testing.T) {
	actual := collectSeq(
		Of(slices.Values([]int{3, 8, 11, 42})).Limit(2),
	)

	expect := []int{3, 8}
	assert.Equal(t, expect, actual)
}

func TestSeqLimitStopsPulling(t *testing.T) {
	var pulled []int
	src := Seq[int](func(yield func(int) bool) {
		for _, v := range []int{3, 8, 11, 42} {
			pulled = append(pulled, v)
			if !yield(v) {
				return
			}
		}
	})

	actual := collectSeq(src.Limit(2))

	expect := []int{3, 8}
	assert.Equal(t, expect, actual)
	assert.Equal(t, expect, pulled, "Limit should stop pulling from its parent once satisfied")
}
