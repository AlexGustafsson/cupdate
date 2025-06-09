package workflow

import (
	"context"
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/models"
)

var _ context.Context = (*Context)(nil)

// Context is a [context.Context] implementation holding additional context for
// a workflow. The value should be seen as being immutable.
type Context struct {
	context.Context

	// Workflow is the current workflow.
	Workflow Workflow
	// WorkflowRun describes the current invocation.
	WorkflowRun models.WorkflowRun
	// Job is the current job.
	Job      Job
	JobIndex int
	// Step is the current step.
	Step      Step
	StepIndex int

	// Outputs holds the outputs of steps and jobs, mapped by their path.
	// Outputs should not be written directly, as they are managed by the workflow
	// runtime.
	// Example:
	//   ctx.Outputs["step.getManifests.manifests"]
	//   ctx.Outputs["jobs.oci.step.getManifests.manifests"]
	Outputs *Map[string, any]

	// Error holds any current error of the context.
	Error error
}

// GetValue returns a value in the ctx.
// If a value does not exist, the type's zero value is returned, with a nil
// error.
func GetValue[T any](ctx Context, name string) (T, error) {
	v, ok := GetAnyValue(ctx, name)

	var ret T
	if !ok {
		return ret, nil
	}

	ret, ok = v.(T)
	if !ok {
		return ret, fmt.Errorf("invalid type %T for value %s of type %T", v, name, ret)
	}

	return ret, nil
}

// GetAnyValue returns a value in the ctx.
func GetAnyValue(ctx Context, name string) (any, bool) {
	v, ok := ctx.Outputs[name]
	return v, ok
}
