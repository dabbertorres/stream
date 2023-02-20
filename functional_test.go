package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	input := map[string]int{
		"foo":   5,
		"bar":   -7,
		"baz":   43,
		"xyzzy": -11,
		"plugh": 512,
	}

	actual := Pipe(
		FromMap(input),
		FilterBy(func(kv KeyValue[string, int]) bool { return kv.Val > 0 }),
		Chain(
			MapBy(ValueMapper[string](func(i int) bool { return i%2 == 0 })),
			AssociateBy(AssociateKeyValue[string, bool]),
		),
	)

	expect := map[string]bool{
		"baz":   false,
		"foo":   false,
		"plugh": true,
	}
	assert.Equal(t, expect, actual)
}
