package workflow

import (
	"errors"
	"log/slog"
	"maps"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Job represents a series of steps.
// The implementation mimics GitHub actions / workflows.
type Job struct {
	// ID is an optional ID by which to reference any outputs created by the job's
	// steps.
	ID string
	// Name is a human-readable name of the job.
	Name string
	// Steps holds all the steps that are run as part of the job.
	Steps []Step
	// DependsOn optionally holds the ids of jobs that must complete for this job
	// to run.
	// Jobs will not run if one of its dependencies fail.
	DependsOn []string
	// If holds a [Condition] that must pass for this job to run.
	If Condition
}

// Run runs the job.
// Returns the modified context and any error.
// Returns [ErrSkipped] if the job was not run.
func (j Job) Run(ctx Context) (Context, error) {
	log := slog.With(slog.String("workflow", ctx.Workflow.Name), slog.String("job", ctx.Job.Name))

	shouldRun := true
	if j.If != nil {
		var err error
		shouldRun, err = testCondition(ctx, j.If)
		if err != nil {
			return ctx, err
		}
	}

	if !shouldRun {
		log.DebugContext(ctx, "Skipped job")
		return ctx, ErrSkipped
	}

	log.DebugContext(ctx, "Running job")

	var jobSpan trace.Span
	ctx.Context, jobSpan = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, otelutil.CupdateWorkflowJobRunSpanName, trace.WithAttributes(otelutil.CupdateWorkflowJobName(j.Name)))
	defer jobSpan.End()

	errs := make([]error, len(j.Steps)*2)

	log.DebugContext(ctx, "Running job steps")
	for i, step := range j.Steps {
		if step.Main == nil {
			continue
		}

		ctx := Context{
			Context: ctx.Context,

			Workflow:    ctx.Workflow,
			WorkflowRun: ctx.WorkflowRun,
			Job:         ctx.Job,
			JobIndex:    ctx.JobIndex,
			Step:        step,
			StepIndex:   i,

			Outputs: ctx.Outputs,

			Error: errors.Join(errs...),
		}

		started := time.Now()
		ctx, err := step.Run(ctx)
		if err == ErrSkipped {
			continue
		}

		ctx.WorkflowRun.Jobs[ctx.JobIndex].Steps[i].Started = started
		ctx.WorkflowRun.Jobs[ctx.JobIndex].Steps[i].DurationSeconds = time.Since(started).Seconds()

		if err != nil {
			errs[i] = err
			ctx.WorkflowRun.Jobs[ctx.JobIndex].Steps[i].Result = models.StepRunResultFailed
			ctx.WorkflowRun.Jobs[ctx.JobIndex].Steps[i].Error = err.Error()
			continue
		}

		ctx.WorkflowRun.Jobs[ctx.JobIndex].Steps[i].Result = models.StepRunResultSucceeded
		outputs = maps.Clone(ctx.Outputs)
	}

	log.DebugContext(ctx, "Running post steps")
	for i, step := range j.Steps {
		if step.Post == nil {
			continue
		}

		ctx := Context{
			Context: ctx.Context,

			Workflow: ctx.Workflow,
			Job:      ctx.Job,
			Step:     step,

			Outputs: outputs,

			Error: errors.Join(errs...),
		}

		// TODO: Represent in workflow run somehow
		err := step.RunPost(ctx)
		if err == ErrSkipped {
			continue
		} else if err != nil {
			errs[len(j.Steps)+i] = err
			continue
		}
	}

	if err := errors.Join(errs...); err != nil {
		jobSpan.SetStatus(codes.Error, "One or more steps failed")
		return ctx, errors.Join(errs...)
	}

	jobSpan.SetStatus(codes.Ok, "")
	ctx.Outputs = outputs
	return ctx, nil
}
