package stream

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	o := Some(37)

	assert.True(t, o.Some())
	assert.False(t, o.None())

	v, ok := o.Get()
	assert.True(t, ok)
	assert.Equal(t, 37, v)
}

func TestNone(t *testing.T) {
	o := None[int]()

	assert.False(t, o.Some())
	assert.True(t, o.None())

	v, ok := o.Get()
	assert.False(t, ok)
	assert.Equal(t, 0, v)
}

func TestOptionalFromPointer(t *testing.T) {
	t.Run("non_nil", func(t *testing.T) {
		val := 37
		o := OptionalFromPointer(&val)

		v, ok := o.Get()
		assert.True(t, ok)
		assert.Equal(t, 37, v)
	})

	t.Run("nil", func(t *testing.T) {
		o := OptionalFromPointer[int](nil)

		assert.True(t, o.None())
	})
}

func TestOptionalMustGet(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		assert.Equal(t, 37, Some(37).MustGet())
	})

	t.Run("none", func(t *testing.T) {
		assert.Panics(t, func() { None[int]().MustGet() })
	})
}

func TestOptionalGetOrDefault(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		assert.Equal(t, 37, Some(37).GetOrDefault(-1))
	})

	t.Run("none", func(t *testing.T) {
		assert.Equal(t, -1, None[int]().GetOrDefault(-1))
	})
}

func TestOptionalGetOrDefaultFunc(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var called bool
		actual := Some(37).GetOrDefaultFunc(func() int { called = true; return -1 })

		assert.Equal(t, 37, actual)
		assert.False(t, called, "default func should not be called when Some")
	})

	t.Run("none", func(t *testing.T) {
		actual := None[int]().GetOrDefaultFunc(func() int { return -1 })

		assert.Equal(t, -1, actual)
	})
}

func TestOptionalIfSome(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var got int
		Some(37).IfSome(func(v int) { got = v })

		assert.Equal(t, 37, got)
	})

	t.Run("none", func(t *testing.T) {
		var called bool
		None[int]().IfSome(func(int) { called = true })

		assert.False(t, called)
	})
}

func TestOptionalIfNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		var called bool
		Some(37).IfNone(func() { called = true })

		assert.False(t, called)
	})

	t.Run("none", func(t *testing.T) {
		var called bool
		None[int]().IfNone(func() { called = true })

		assert.True(t, called)
	})
}

func TestOptionalMarshalJSON(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		data, err := json.Marshal(Some(37))
		assert.NoError(t, err)
		assert.JSONEq(t, `37`, string(data))
	})

	t.Run("none", func(t *testing.T) {
		data, err := json.Marshal(None[int]())
		assert.NoError(t, err)
		assert.JSONEq(t, `null`, string(data))
	})

	t.Run("struct_field", func(t *testing.T) {
		type T struct {
			Foo Optional[int] `json:"foo"`
		}

		data, err := json.Marshal(T{Foo: Some(37)})
		assert.NoError(t, err)
		assert.JSONEq(t, `{"foo": 37}`, string(data))

		data, err = json.Marshal(T{Foo: None[int]()})
		assert.NoError(t, err)
		assert.JSONEq(t, `{"foo": null}`, string(data))
	})
}

func TestOptionalUnmarshalJSON(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		var o Optional[int]
		err := json.Unmarshal([]byte(`37`), &o)

		assert.NoError(t, err)
		v, ok := o.Get()
		assert.True(t, ok)
		assert.Equal(t, 37, v)
	})

	t.Run("null", func(t *testing.T) {
		o := Some(37)
		err := json.Unmarshal([]byte(`null`), &o)

		assert.NoError(t, err)
		assert.True(t, o.None())
	})

	t.Run("invalid", func(t *testing.T) {
		var o Optional[int]
		err := json.Unmarshal([]byte(`"not a number"`), &o)

		assert.Error(t, err)
	})
}
