package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqFilter(t *testing.T) {
	actual := collectSeq(
		Of(slices.Values([]int{3, 8, 11, 42})).Filter(func(i int) bool { return i%2 == 0 }),
	)

	expect := []int{8, 42}
	assert.Equal(t, expect, actual)
}
