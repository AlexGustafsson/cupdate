package workflow

import (
	"context"
	"fmt"
	"strings"
)

var _ context.Context = (*Context)(nil)

type Context struct {
	context.Context

	// Workflow is the current workflow.
	Workflow Workflow
	// Job is the current job.
	Job Job
	// Step is the current step.
	Step Step

	// Outputs holds the outputs of steps and jobs, mapped by their path.
	// Outputs should not be written directly, as they are managed by the workflow
	// runtime.
	// Example:
	//   ctx.Outputs["step.getManifests.manifests"]
	//   ctx.Outputs["jobs.oci.step.getManifests.manifests"]
	Outputs map[string]any

	// Values holds values (variables) stored by calling Store on a step to store
	// a named output. Values can later be used as inputs.
	// Values should not be written directly, as they are managed by the workflow
	// runtime.
	// Example:
	//    FetchTitlesFromIMDB().Store("titles", "movieTitles")
	//    CreateReportFromTitles("movieTitles")
	Values map[string]any

	// Error holds any current error of the context.
	Error error
}

// GetValue returns a value or output in the ctx.
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
		return ret, fmt.Errorf("invalid value")
	}

	return ret, nil
}

// GetAnyValue returns a value or output in the ctx.
func GetAnyValue(ctx Context, name string) (any, bool) {
	var v any
	var ok bool
	if strings.Contains(name, ".") {
		v, ok = ctx.Outputs[name]
	} else {
		v, ok = ctx.Values[name]
	}

	return v, ok
}
