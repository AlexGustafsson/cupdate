package workflow

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepWithID(t *testing.T) {
	assert.Equal(t, Step{ID: "1"}, Step{}.WithID("1"))
}

func TestStepWithCondition(t *testing.T) {
	assert.Panics(t, func() {
		Step{
			If: "already set",
		}.WithCondition("already set")
	})

	assert.Equal(t, Step{If: "step.1.value"}, Step{}.WithCondition("step.1.value"))
}

func TestStepWith(t *testing.T) {
	assert.Equal(t, Step{Inputs: map[string]any{"foo": "bar"}}, Step{}.With("foo", "bar"))
	assert.Equal(t, Step{Inputs: map[string]any{"foo": "baz"}}, Step{Inputs: map[string]any{"foo": "bar"}}.With("foo", "baz"))
}

func TestStepRun(t *testing.T) {
	testCases := []struct {
		Name               string
		Step               Step
		Context            Context
		AssertExpectations func(*testing.T, Context)
		Error              error
	}{
		{
			Name: "Steps should not run if previous step failed",
			Context: Context{
				Error: fmt.Errorf("previous step failed"),
			},
			Error: ErrSkipped,
		},
		{
			Name: "Skip based on conditional",
			Step: Step{
				If: ConditionFunc(func(ctx Context) (bool, error) {
					return false, nil
				}),
			},
			Context: Context{},
			Error:   ErrSkipped,
		},
		{
			Name: "Conditional may fail",
			Step: Step{
				If: ConditionFunc(func(ctx Context) (bool, error) {
					return false, fmt.Errorf("conditional failed")
				}),
			},
			Context: Context{},
			Error:   fmt.Errorf("conditional failed"),
		},
		{
			Name: "Main may fail",
			Step: Step{
				Main: func(ctx Context) (Command, error) {
					return nil, fmt.Errorf("main failed")
				},
			},
			Context: Context{
				Context: context.TODO(),
			},
			Error: fmt.Errorf("main failed"),
		},
		{
			Name: "Commands are run",
			Step: Step{
				Main: func(ctx Context) (Command, error) {
					return SetOutput("output", 42), nil
				},
			},
			Context: Context{
				Context: context.TODO(),
				Step: Step{
					ID: "1",
				},
				Outputs: make(map[string]any),
			},
			AssertExpectations: func(t *testing.T, ctx Context) {
				assert.Equal(t, map[string]any{"step.1.output": 42}, ctx.Outputs)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual, err := testCase.Step.Run(testCase.Context)
			assert.Equal(t, err, testCase.Error)
			if testCase.AssertExpectations != nil {
				testCase.AssertExpectations(t, actual)
			}
		})
	}
}

func TestStepRunPost(t *testing.T) {
	testCases := []struct {
		Name    string
		Step    Step
		Context Context
		Error   error
	}{
		{
			Name: "Happy path",
			Step: Step{
				Post: func(ctx Context) error {
					return nil
				},
			},
			Context: Context{
				Context: context.TODO(),
			},
			Error: nil,
		},
		{
			Name: "Post steps should not run if previous step failed",
			Context: Context{
				Error: fmt.Errorf("previous step failed"),
			},
			Error: ErrSkipped,
		},
		{
			Name: "Skip based on conditional",
			Step: Step{
				PostIf: ConditionFunc(func(ctx Context) (bool, error) {
					return false, nil
				}),
			},
			Context: Context{},
			Error:   ErrSkipped,
		},
		{
			Name: "Conditional may fail",
			Step: Step{
				PostIf: ConditionFunc(func(ctx Context) (bool, error) {
					return false, fmt.Errorf("conditional failed")
				}),
			},
			Context: Context{},
			Error:   fmt.Errorf("conditional failed"),
		},
		{
			Name: "Post may fail",
			Step: Step{
				Post: func(ctx Context) error {
					return fmt.Errorf("main failed")
				},
			},
			Context: Context{
				Context: context.TODO(),
			},
			Error: fmt.Errorf("main failed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := testCase.Step.RunPost(testCase.Context)
			assert.Equal(t, err, testCase.Error)
		})
	}
}

func TestStepFunc(t *testing.T) {
	f := StepFunc(func(ctx Context) (Command, error) {
		return nil, nil
	})

	// NOTE: We cannot assert functions in Go
	assert.NotNil(t, Run(f).Main)
}
