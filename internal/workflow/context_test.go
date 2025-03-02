package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	ctx := Context{
		Outputs: map[string]any{
			"foo": "bar",
			"bar": 1,
			"baz": (*string)(nil),
		},
	}

	stringValue, err := GetValue[string](ctx, "foo")
	assert.Equal(t, "bar", stringValue)
	assert.NoError(t, err)

	intValue, err := GetValue[int](ctx, "bar")
	assert.Equal(t, 1, intValue)
	assert.NoError(t, err)

	nilValue, err := GetValue[*string](ctx, "baz")
	assert.Equal(t, (*string)(nil), nilValue)
	assert.NoError(t, err)

	notFound, err := GetValue[string](ctx, "not-found")
	assert.Equal(t, "", notFound)
	assert.NoError(t, err)

	wrongType, err := GetValue[int](ctx, "foo")
	assert.Equal(t, 0, wrongType)
	assert.Error(t, err)
}

func TestGetAnyValue(t *testing.T) {
	ctx := Context{
		Outputs: map[string]any{
			"foo": "bar",
		},
	}

	value, ok := GetAnyValue(ctx, "foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)

	value, ok = GetAnyValue(ctx, "bar")
	assert.Nil(t, value)
	assert.False(t, ok)
}
