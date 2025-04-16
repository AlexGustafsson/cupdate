package workflow

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestJobRun(t *testing.T) {
	testCases := []struct {
		Name               string
		Job                Job
		Context            Context
		AssertExpectations func(*testing.T, Context)
		Error              error
	}{
		{
			Name: "Happy path",
			Job: Job{
				Steps: []Step{
					{
						ID: "1",
						Main: func(ctx Context) (Command, error) {
							return SetOutput("foo", "bar"), nil
						},
					},
					{
						ID: "2",
						Main: func(ctx Context) (Command, error) {
							return SetOutput("foo", "bar"), nil
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
								{},
							},
						},
					},
				},
			},
			AssertExpectations: func(t *testing.T, ctx Context) {
				assert.Equal(t, models.StepRunResultSucceeded, ctx.WorkflowRun.Jobs[0].Steps[0].Result)
				assert.Equal(t, models.StepRunResultSucceeded, ctx.WorkflowRun.Jobs[0].Steps[1].Result)
				assert.Equal(t, map[string]any{"step.1.foo": "bar", "step.2.foo": "bar"}, ctx.Outputs)
			},
		},
		{
			Name: "Skip based on conditional",
			Job: Job{
				If: ConditionFunc(func(ctx Context) (bool, error) {
					return false, nil
				}),
			},
			Context: Context{},
			Error:   ErrSkipped,
		},
		{
			Name: "Conditional may fail",
			Job: Job{
				If: ConditionFunc(func(ctx Context) (bool, error) {
					return false, fmt.Errorf("conditional failed")
				}),
			},
			Context: Context{},
			Error:   fmt.Errorf("conditional failed"),
		},
		{
			Name: "Skip steps without main",
			Job: Job{
				Steps: []Step{
					{
						Main: nil,
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
							},
						},
					},
				},
			},
		},
		{
			Name: "Skipped steps are marked as such",
			Job: Job{
				Steps: []Step{
					{
						If: Condition(func(ctx Context) (bool, error) {
							return false, nil
						}),
						Main: func(ctx Context) (Command, error) {
							t.Fail()
							return nil, nil
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{
									// Setting it as the default is how the workflow does it
									Result: models.StepRunResultSkipped,
								},
							},
						},
					},
				},
			},
			AssertExpectations: func(t *testing.T, ctx Context) {
				assert.Equal(t, models.StepRunResultSkipped, ctx.WorkflowRun.Jobs[0].Steps[0].Result)
			},
		},
		{
			Name: "Step errors are reported",
			Job: Job{
				Steps: []Step{
					{
						Main: func(ctx Context) (Command, error) {
							return nil, fmt.Errorf("failed")
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
							},
						},
					},
				},
			},
			AssertExpectations: func(t *testing.T, ctx Context) {
				assert.Equal(t, models.StepRunResultFailed, ctx.WorkflowRun.Jobs[0].Steps[0].Result)
				assert.Equal(t, "failed", ctx.WorkflowRun.Jobs[0].Steps[0].Error)
			},
			Error: errors.Join(fmt.Errorf("failed")),
		},
		// TODO: This test doesn't clearly assert anything - it does trigger paths
		// in the coverage though, and implicitly tests that they work
		{
			Name: "Step post runs are run",
			Job: Job{
				Steps: []Step{
					{
						Main: func(ctx Context) (Command, error) {
							return nil, nil
						},
						Post: func(ctx Context) error {
							return nil
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
							},
						},
					},
				},
			},
			Error: nil,
		},
		// TODO: This test doesn't clearly assert anything - it does trigger paths
		// in the coverage though, and returns an error, implicitly asserting that
		// it is handled as expected
		{
			Name: "Step post runs can be skipped",
			Job: Job{
				Steps: []Step{
					{
						Main: func(ctx Context) (Command, error) {
							return nil, nil
						},
						Post: func(ctx Context) error {
							return ErrSkipped
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
							},
						},
					},
				},
			},
			Error: nil,
		},
		{
			Name: "Step post runs can fail",
			Job: Job{
				Steps: []Step{
					{
						Main: func(ctx Context) (Command, error) {
							return nil, nil
						},
						Post: func(ctx Context) error {
							return fmt.Errorf("failed")
						},
					},
				},
			},
			Context: Context{
				Context: context.TODO(),
				Outputs: make(map[string]any),
				WorkflowRun: models.WorkflowRun{
					Jobs: []models.JobRun{
						{
							Steps: []models.StepRun{
								{},
							},
						},
					},
				},
			},
			AssertExpectations: func(t *testing.T, ctx Context) {
				// TODO: Post run failures are currently not reported
				assert.Equal(t, models.StepRunResultSucceeded, ctx.WorkflowRun.Jobs[0].Steps[0].Result)
				assert.Equal(t, "", ctx.WorkflowRun.Jobs[0].Steps[0].Error)
			},
			Error: errors.Join(fmt.Errorf("failed")),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual, err := testCase.Job.Run(testCase.Context)
			assert.Equal(t, testCase.Error, err)
			if testCase.AssertExpectations != nil {
				testCase.AssertExpectations(t, actual)
			}
		})
	}
}
