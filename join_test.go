package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	actual := collectSeq(
		Join(
			OfSlice([]int{7, 11, 19}),
			OfSlice([]int{42, 23, 9}),
			OfSlice([]int{3}),
		),
	)

	expect := []int{7, 11, 19, 42, 23, 9, 3}
	assert.Equal(t, expect, actual)
}
