package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	calls := ""

	command := Batch(
		func(ctx Context) {
			calls += "a"
		},
		func(ctx Context) {
			calls += "b"
		},
		func(ctx Context) {
			calls += "c"
		},
	)

	command(Context{})

	assert.Equal(t, "abc", calls)
}

func TestSetOutput(t *testing.T) {
	ctx := Context{
		Step: Step{
			ID: "1",
		},
		Outputs: make(map[string]any),
	}

	command := SetOutput("foo", "bar")

	command(ctx)

	assert.Equal(t, map[string]any{"step.1.foo": "bar"}, ctx.Outputs)
}
