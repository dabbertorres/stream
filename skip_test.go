package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqSkip(t *testing.T) {
	actual := collectSeq(
		Of(slices.Values([]int{3, 8, 11, 42})).Skip(2),
	)

	expect := []int{11, 42}
	assert.Equal(t, expect, actual)
}
