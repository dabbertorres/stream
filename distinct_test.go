package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinct(t *testing.T) {
	actual := collectSeq(Distinct(OfSlice([]int{3, 8, 3, 11, 8, 8, 42})))

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}
