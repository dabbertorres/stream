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
		ApplyRight(
			Stream[KeyValue[string, int]].Filter,
			func(kv KeyValue[string, int]) bool { return kv.Val > 0 },
		),
		Chain(
			Chain(
				Map[KeyValue[string, bool], KeyValue[string, int]],
				ApplyRight(
					Mapper[KeyValue[string, int], KeyValue[string, bool]].By,
					ValueMapper[bool, string](func(i int) bool { return i%2 == 0 }),
				),
			),
			Chain(
				Associate[string, bool, KeyValue[string, bool]],
				ApplyRight(
					Associater[KeyValue[string, bool], string, bool].By,
					AssociateKeyValue[string, bool],
				),
			),
		),
	)

	expect := map[string]bool{
		"baz":   false,
		"foo":   false,
		"plugh": true,
	}
	assert.Equal(t, expect, actual)
}
