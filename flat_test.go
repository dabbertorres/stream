package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqFlatMap(t *testing.T) {
	actual := collectSeq(
		OfSlice([]int{5, 10, 2, 90}).FlatMap(func(x int) Seq[bool] {
			return OfSlice([]bool{x%2 == 0, x > 5})
		}),
	)

	expect := []bool{false, false, true, true, true, false, true, true}
	assert.Equal(t, expect, actual)
}
