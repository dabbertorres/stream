package stream

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOfSlice(t *testing.T) {
	actual := collectSeq(OfSlice([]int{3, 8, 11, 42}))

	expect := []int{3, 8, 11, 42}
	assert.Equal(t, expect, actual)
}

func TestOfMap(t *testing.T) {
	actual := collectSeq(OfMap(map[string]int{"foo": 1}))

	expect := []Pair[string, int]{{Key: "foo", Val: 1}}
	assert.Equal(t, expect, actual)
}

func TestOfOptional(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		actual := collectSeq(OfOptional(Some(37)))

		expect := []int{37}
		assert.Equal(t, expect, actual)
	})

	t.Run("none", func(t *testing.T) {
		actual := collectSeq(OfOptional(None[int]()))

		assert.Empty(t, actual)
	})
}

func TestOfFunc(t *testing.T) {
	i := 1
	actual := collectSeq(OfFunc(func() Optional[int] {
		v := i
		i *= 2
		if v <= 64 {
			return Some(v)
		}
		return None[int]()
	}))

	expect := []int{1, 2, 4, 8, 16, 32, 64}
	assert.Equal(t, expect, actual)
}

func TestOfDecoder(t *testing.T) {
	type bar struct {
		Foo int `json:"foo"`
	}

	data := []byte(`
		{"foo": 7}
		{"foo": 5}
		{"foo": -13}
	`)

	var onErrCalled bool
	actual := collectSeq(OfDecoder[bar](json.NewDecoder(bytes.NewReader(data)), func(error) { onErrCalled = true }))

	expect := []bar{{Foo: 7}, {Foo: 5}, {Foo: -13}}
	assert.Equal(t, expect, actual)
	assert.False(t, onErrCalled)
}

func TestOfDecoderError(t *testing.T) {
	type bar struct {
		Foo int `json:"foo"`
	}

	data := []byte(`{"foo": 7} not valid json`)

	var gotErr error
	actual := collectSeq(OfDecoder[bar](json.NewDecoder(bytes.NewReader(data)), func(err error) { gotErr = err }))

	expect := []bar{{Foo: 7}}
	assert.Equal(t, expect, actual)
	assert.Error(t, gotErr)
}
