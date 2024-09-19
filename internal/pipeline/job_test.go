package pipeline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct {
	FieldA string
	FieldB int
}

func TestContextGetOutput(t *testing.T) {
	ctx := newContext[any](context.TODO(), nil, nil)

	// positive: string
	stringValue := "foo bar"
	var stringReturnValue string
	ctx.SetOutput("value", stringValue)
	assert.True(t, ctx.GetOutput("value", &stringReturnValue))
	assert.Equal(t, stringValue, stringReturnValue)

	// positive: int
	intValue := 1
	var intReturnValue int
	ctx.SetOutput("value", intValue)
	assert.True(t, ctx.GetOutput("value", &intReturnValue))
	assert.Equal(t, intValue, intReturnValue)

	// positive: struct
	structValue := testType{
		FieldA: "a",
		FieldB: 2,
	}
	var structReturnValue testType
	ctx.SetOutput("value", structValue)
	assert.True(t, ctx.GetOutput("value", &structReturnValue))
	assert.Equal(t, structValue, structReturnValue)

	// positive: struct pointer
	structPointerValue := &testType{
		FieldA: "a",
		FieldB: 2,
	}
	var structPointerReturnValue *testType
	ctx.SetOutput("value", structPointerValue)
	assert.True(t, ctx.GetOutput("value", &structPointerReturnValue))
	assert.Equal(t, structPointerValue, structPointerReturnValue)

	// negative: int to struct
	ctx.SetOutput("value", "not a struct")
	assert.False(t, ctx.GetOutput("value", &structReturnValue))
}
