package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqReduce(t *testing.T) {
	actual := Of(slices.Values([]int{3, 8, 11, 42})).
		Reduce(0, func(cumulative, next int) int { return cumulative + next })

	assert.Equal(t, 64, actual)
}
