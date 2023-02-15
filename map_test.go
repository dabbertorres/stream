package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromMap(t *testing.T) {
	input := map[string]int{
		"foo":   5,
		"bar":   -7,
		"baz":   43,
		"xyzzy": -11,
		"plugh": 513,
	}

	actual := FromMap(input).
		Sorted(KeyValueLess[string, int]).
		Collect()

	expect := []KeyValue[string, int]{
		{"bar", -7},
		{"baz", 43},
		{"foo", 5},
		{"plugh", 513},
		{"xyzzy", -11},
	}

	assert.Equal(t, expect, actual)
}
