package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInput(t *testing.T) {
	// Err is propagated
	errValue, err := GetInput[string](Context{}, "foo", true)
	assert.Equal(t, "", errValue)
	assert.Error(t, err)

	// Nil value
	nilValue, err := GetInput[*string](Context{Step: Step{Inputs: map[string]Input{"foo": nil}}}, "foo", true)
	assert.Equal(t, (*string)(nil), nilValue)
	assert.NoError(t, err)

	// Invalid type
	invalidType, err := GetInput[int](Context{Step: Step{Inputs: map[string]Input{"foo": "bar"}}}, "foo", true)
	assert.Equal(t, 0, invalidType)
	assert.Error(t, err)

	// Valid
	stringValue, err := GetInput[string](Context{Step: Step{Inputs: map[string]Input{"foo": "bar"}}}, "foo", true)
	assert.Equal(t, "bar", stringValue)
	assert.NoError(t, err)
}

func TestGetAnyInput(t *testing.T) {
	// No inputs for step, value required
	value, err := GetAnyInput(Context{}, "foo", true)
	assert.Nil(t, value)
	assert.Error(t, err)

	// No inputs for step, value not required
	value, err = GetAnyInput(Context{}, "foo", false)
	assert.Nil(t, value)
	assert.NoError(t, err)

	// Step has inputs, but value not found whilst required
	value, err = GetAnyInput(Context{Step: Step{Inputs: map[string]Input{}}}, "foo", true)
	assert.Nil(t, value)
	assert.Error(t, err)

	// Step has inputs, but value not found whilst not required
	value, err = GetAnyInput(Context{Step: Step{Inputs: map[string]Input{}}}, "foo", false)
	assert.Nil(t, value)
	assert.NoError(t, err)

	// Value is verbatim
	value, err = GetAnyInput(Context{Step: Step{Inputs: map[string]Input{"foo": "bar"}}}, "foo", false)
	assert.Equal(t, "bar", value)
	assert.NoError(t, err)

	// Value is existing ref
	value, err = GetAnyInput(Context{Outputs: map[string]any{"step.1.foo": "bar"}, Step: Step{Inputs: map[string]Input{"foo": Ref{"step.1.foo"}}}}, "foo", false)
	assert.Equal(t, "bar", value)
	assert.NoError(t, err)

	// Value is missing ref
	value, err = GetAnyInput(Context{Outputs: map[string]any{}, Step: Step{Inputs: map[string]Input{"foo": Ref{"step.1.foo"}}}}, "foo", false)
	assert.Equal(t, nil, value)
	assert.NoError(t, err)
}
