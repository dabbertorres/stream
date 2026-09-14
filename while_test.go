package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqDropWhile(t *testing.T) {
	actual := collectSeq(OfSlice([]int{3, 8, 11, 42}).DropWhile(func(i int) bool { return i < 10 }))

	expect := []int{11, 42}
	assert.Equal(t, expect, actual)
}

func TestSeqTakeWhile(t *testing.T) {
	actual := collectSeq(OfSlice([]int{3, 8, 11, 42}).TakeWhile(func(i int) bool { return i < 10 }))

	expect := []int{3, 8}
	assert.Equal(t, expect, actual)
}

func TestSeqTakeWhileStopsPulling(t *testing.T) {
	var pulled []int
	src := Seq[int](func(yield func(int) bool) {
		for _, v := range []int{3, 8, 11, 42} {
			pulled = append(pulled, v)
			if !yield(v) {
				return
			}
		}
	})

	actual := collectSeq(src.TakeWhile(func(i int) bool { return i < 10 }))

	assert.Equal(t, []int{3, 8}, actual)
	assert.Equal(t, []int{3, 8, 11}, pulled, "TakeWhile should stop pulling once the predicate fails")
}
