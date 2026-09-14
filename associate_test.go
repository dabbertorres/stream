package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqAssociate(t *testing.T) {
	type T struct {
		Key string
		Val int
	}

	actual := collectSeq(
		Of(slices.Values([]T{
			{Key: "foo", Val: 37},
			{Key: "bar", Val: -10},
		})).Associate(func(elem T) (string, int) { return elem.Key, elem.Val }),
	)

	expect := []Pair[string, int]{
		{Key: "foo", Val: 37},
		{Key: "bar", Val: -10},
	}
	assert.Equal(t, expect, actual)
}

func TestToMap(t *testing.T) {
	type T struct {
		Key string
		Val int
	}

	actual := ToMap(
		OfSlice([]T{
			{Key: "foo", Val: 37},
			{Key: "bar", Val: -10},
		}).Associate(func(elem T) (string, int) { return elem.Key, elem.Val }),
	)

	expect := map[string]int{
		"foo": 37,
		"bar": -10,
	}
	assert.Equal(t, expect, actual)
}

func TestToMapLastWriteWins(t *testing.T) {
	actual := ToMap(OfSlice([]Pair[string, int]{
		{Key: "foo", Val: 1},
		{Key: "foo", Val: 2},
	}))

	expect := map[string]int{"foo": 2}
	assert.Equal(t, expect, actual)
}
