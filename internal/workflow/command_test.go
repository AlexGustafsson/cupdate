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
	testCases := []struct {
		Name            string
		StepID          string
		Command         Command
		ExpectedOutputs map[string]any
	}{
		{
			Name:            "Happy path",
			StepID:          "1",
			Command:         SetOutput("foo", "bar"),
			ExpectedOutputs: map[string]any{"step.1.foo": "bar"},
		},
		{
			Name:            "Set value, no step id",
			StepID:          "",
			Command:         SetOutput("foo", "bar"),
			ExpectedOutputs: map[string]any{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx := Context{
				Step: Step{
					ID: testCase.StepID,
				},
				Outputs: make(map[string]any),
			}

			testCase.Command(ctx)

			assert.Equal(t, testCase.ExpectedOutputs, ctx.Outputs)
		})
	}
}
