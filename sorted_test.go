package stream

import (
	"cmp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqSorted(t *testing.T) {
	t.Run("is_sorted", func(t *testing.T) {
		actual := collectSeq(OfSlice([]int{3, 8, 11, 42}).Sorted(cmp.Less))

		expect := []int{3, 8, 11, 42}
		assert.Equal(t, expect, actual)
	})

	t.Run("not_sorted", func(t *testing.T) {
		actual := collectSeq(OfSlice([]int{42, 3, 11, 8}).Sorted(cmp.Less))

		expect := []int{3, 8, 11, 42}
		assert.Equal(t, expect, actual)
	})
}
