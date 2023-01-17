package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssociate(t *testing.T) {
	type T struct {
		Key string
		Val int
	}

	actual := Associate(
		FromSlice([]T{{Key: "foo", Val: 37}, {Key: "bar", Val: -10}, {Key: "baz", Val: 45}}),
		func(elem T) (string, int) { return elem.Key, elem.Val },
	)

	expect := map[string]int{
		"foo": 37,
		"bar": -10,
		"baz": 45,
	}
	assert.Equal(t, expect, actual)
}
