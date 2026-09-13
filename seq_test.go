package stream

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func collectSeq[T any](s Seq[T]) []T {
	var out []T
	for v := range s {
		out = append(out, v)
	}
	return out
}

func TestOf(t *testing.T) {
	actual := collectSeq(Of(slices.Values([]int{3, 8, 11, 42})))

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestOfChan(t *testing.T) {
	ch := make(chan int)
	go func() {
		defer close(ch)
		ch <- 3
		ch <- 8
		ch <- 11
	}()

	actual := collectSeq(OfChan[<-chan int](ch))

	expect := []int{3, 8, 11}
	assert.Equal(t, expect, actual)
}

func TestOfKeyed(t *testing.T) {
	src := map[string]int{"foo": 1}

	actual := collectSeq(OfKeyed(maps.All(src)))

	expect := []Pair[string, int]{{Key: "foo", Val: 1}}
	assert.Equal(t, expect, actual)
}

func TestSeqChained(t *testing.T) {
	actual := collectSeq(
		Of(slices.Values([]int{3, 8, 11, 24, 37, 42})).
			Filter(func(i int) bool { return i >= 10 }).
			Skip(1).
			Limit(2).
			Map(func(i int) int { return i * 2 }),
	)

	expect := []int{48, 74}
	assert.Equal(t, expect, actual)
}
