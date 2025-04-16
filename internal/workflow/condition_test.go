package workflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlways(t *testing.T) {
	shouldRun, err := Always(Context{})
	assert.True(t, shouldRun)
	assert.NoError(t, err)
}

func TestValueExists(t *testing.T) {
	ctx := Context{
		Outputs: map[string]any{
			"foo": "bar",
		},
	}

	condition := ValueExists("foo")
	shouldRun, err := condition(ctx)
	assert.True(t, shouldRun)
	assert.NoError(t, err)

	condition = ValueExists("bar")
	shouldRun, err = condition(ctx)
	assert.False(t, shouldRun)
	assert.NoError(t, err)

	// Job has already failed
	ctx.Error = fmt.Errorf("some error")
	condition = ValueExists("foo")
	shouldRun, err = condition(ctx)
	assert.False(t, shouldRun)
	assert.NoError(t, err)
}

func TestTestCondition(t *testing.T) {
	testCases := []struct {
		Name          string
		Condition     Condition
		Context       Context
		Expected      bool
		ExpectedError error
	}{
		{
			Name: "Condition func is called",
			Condition: ConditionFunc(func(ctx Context) (bool, error) {
				assert.Equal(t, "1", ctx.Step.ID)
				return true, nil
			}),
			Context: Context{
				Step: Step{
					ID: "1",
				},
			},
			Expected:      true,
			ExpectedError: nil,
		},
		{
			Name: "Condition func may fail",
			Condition: ConditionFunc(func(ctx Context) (bool, error) {
				return false, fmt.Errorf("failed")
			}),
			Context:       Context{},
			Expected:      false,
			ExpectedError: fmt.Errorf("failed"),
		},
		{
			Name:      "Valid bool reference",
			Condition: "steps.1.cond",
			Context: Context{
				Outputs: map[string]any{"steps.1.cond": true},
			},
			Expected:      true,
			ExpectedError: nil,
		},
		{
			Name:      "Invalid bool reference",
			Condition: "steps.1.cond",
			Context: Context{
				Outputs: map[string]any{"steps.1.cond": "not a bool"},
			},
			Expected:      false,
			ExpectedError: fmt.Errorf("invalid type string for value steps.1.cond of type bool"),
		},
		{
			Name:          "Unsupported condition",
			Condition:     42,
			Context:       Context{},
			Expected:      false,
			ExpectedError: fmt.Errorf("invalid condition type int expected string or func(workflow.Context) (bool, error)"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual, err := testCondition(testCase.Context, testCase.Condition)
			assert.Equal(t, testCase.Expected, actual)
			assert.Equal(t, testCase.ExpectedError, err)
		})
	}
}
