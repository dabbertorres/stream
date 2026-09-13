//go:build !noseqdistinctunsafe

package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqDistinct(t *testing.T) {
	actual := collectSeq(OfSlice([]int{3, 8, 3, 11, 8, 8, 42}).Distinct())

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}
