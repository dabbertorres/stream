package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqMap(t *testing.T) {
	actual := collectSeq(
		Of(slices.Values([]int{5, 10, 3, 90})).Map(func(x int) bool { return x%2 == 0 }),
	)

	expect := []bool{false, true, false, true}
	assert.Equal(t, expect, actual)
}
