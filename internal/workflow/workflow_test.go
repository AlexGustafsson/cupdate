package workflow

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var _ sdktrace.IDGenerator = (*MockTraceIDGenerator)(nil)

type MockTraceIDGenerator struct {
	mock.Mock
}

// NewIDs implements trace.IDGenerator.
func (m *MockTraceIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	args := m.Called(ctx)
	return args.Get(0).(trace.TraceID), args.Get(1).(trace.SpanID)
}

// NewSpanID implements trace.IDGenerator.
func (m *MockTraceIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	args := m.Called(ctx, traceID)
	return args.Get(0).(trace.SpanID)
}

func TestWorkflowRun(t *testing.T) {
	testCases := []struct {
		Name        string
		Workflow    Workflow
		ExpectedRun models.WorkflowRun
		Error       error
	}{
		{
			Name: "Happy path",
			Workflow: Workflow{
				Jobs: []Job{
					{
						ID: "job1",
						Steps: []Step{
							{
								ID: "step1",
								Main: func(ctx Context) (Command, error) {
									return SetOutput("foo", "bar"), nil
								},
							},
						},
					},
					{
						ID: "job2",
						Steps: []Step{
							{
								ID:     "step2",
								Inputs: map[string]Input{"foo": Ref{Key: "job.job1.step.step1.foo"}},
								Main: func(ctx Context) (Command, error) {
									value, err := GetInput[string](ctx, "foo", true)
									assert.Equal(t, "bar", value)
									assert.NoError(t, err)
									return nil, nil
								},
							},
						},
						DependsOn: []string{"job1"},
					},
					{
						ID: "job3",
						If: ConditionFunc(func(ctx Context) (bool, error) {
							return false, nil
						}),
						Steps: []Step{
							{},
						},
					},
					{
						ID: "job4",
						Steps: []Step{
							{},
						},
						DependsOn: []string{"job3"},
					},
				},
			},
			ExpectedRun: models.WorkflowRun{
				Result: models.WorkflowRunResultSucceeded,
				Jobs: []models.JobRun{
					{
						JobID:  "job1",
						Result: models.JobRunResultSucceeded,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSucceeded,
							},
						},
						DependsOn: []string{},
					},
					{
						JobID:  "job2",
						Result: models.JobRunResultSucceeded,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSucceeded,
							},
						},
						DependsOn: []string{"job1"},
					},
					{
						JobID:  "job3",
						Result: models.JobRunResultSkipped,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSkipped,
							},
						},
						DependsOn: []string{},
					},
					{
						JobID:  "job4",
						Result: models.JobRunResultSkipped,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSkipped,
							},
						},
						DependsOn: []string{"job3"},
					},
				},
			},
		},
		{
			Name: "Step fails if context deadline is exceeded",
			Workflow: Workflow{
				Jobs: []Job{
					{
						ID: "job1",
						Steps: []Step{
							{
								ID: "step1",
								Main: func(ctx Context) (Command, error) {
									// Wait for context to be canceled, but don't return an error
									// to make sure the second job fails whilst waiting for this
									<-ctx.Done()
									return nil, nil
								},
							},
						},
					},
					{
						ID: "job2",
						Steps: []Step{
							{
								ID: "step1",
								Main: func(ctx Context) (Command, error) {
									return nil, nil
								},
							},
						},
						DependsOn: []string{"job1"},
					},
				},
			},
			ExpectedRun: models.WorkflowRun{
				Result: models.WorkflowRunResultFailed,
				Jobs: []models.JobRun{
					{
						JobID:  "job1",
						Result: models.JobRunResultSucceeded,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSucceeded,
							},
						},
						DependsOn: []string{},
					},
					{
						JobID:  "job2",
						Result: models.JobRunResultFailed,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultSkipped,
								Error:  "",
							},
						},
						DependsOn: []string{"job1"},
					},
				},
			},
			Error: errors.Join(context.DeadlineExceeded),
		},
		{
			Name: "Failing step fails job",
			Workflow: Workflow{
				Jobs: []Job{
					{
						ID: "job1",
						Steps: []Step{
							{
								ID: "step1",
								Main: func(ctx Context) (Command, error) {
									return nil, fmt.Errorf("failed")
								},
							},
						},
					},
				},
			},
			ExpectedRun: models.WorkflowRun{
				Result: models.WorkflowRunResultFailed,
				Jobs: []models.JobRun{
					{
						JobID:  "job1",
						Result: models.JobRunResultFailed,
						Steps: []models.StepRun{
							{
								Result: models.StepRunResultFailed,
								Error:  "failed",
							},
						},
						DependsOn: []string{},
					},
				},
			},
			Error: errors.Join(errors.Join(fmt.Errorf("failed"))),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
			actual, err := testCase.Workflow.Run(ctx)
			cancel()

			// Remove times to make them easily comparable
			actual.Started = time.Time{}
			actual.DurationSeconds = 0
			for i := range actual.Jobs {
				actual.Jobs[i].Started = time.Time{}
				actual.Jobs[i].DurationSeconds = 0
				for j := range actual.Jobs[i].Steps {
					actual.Jobs[i].Steps[j].Started = time.Time{}
					actual.Jobs[i].Steps[j].DurationSeconds = 0
				}
			}

			assert.Equal(t, testCase.ExpectedRun, actual)
			assert.Equal(t, testCase.Error, err)
		})
	}
}
