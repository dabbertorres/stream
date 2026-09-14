//go:build !noseqdistinctunsafe

package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqDistinctUnsafe(t *testing.T) {
	actual := collectSeq(OfSlice([]int{3, 8, 3, 11, 8, 8, 42}).DistinctUnsafe())

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}
