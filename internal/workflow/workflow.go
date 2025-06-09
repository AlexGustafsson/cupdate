package workflow

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var ErrDependentJobFailed = errors.New("dependent job failed")

// Workflow is a generic way to represent running tasks with or without
// dependencies. The implementation mimics GitHub actions / workflows.
type Workflow struct {
	// Name is the human-readable name of the workflow.
	Name string
	// Jobs holds all the jobs that the workflow should run.
	// Jobs are started in the order defined by their dependencies.
	Jobs []Job
}

// Run executes a workflow and (always) returns a run description and any error
// that caused the workflow to fail.
func (w Workflow) Run(ctx context.Context) (models.WorkflowRun, error) {
	ctx, span := otel.Tracer(otelutil.DefaultScope).Start(ctx, otelutil.CupdateWorkflowRunSpanName, trace.WithAttributes(otelutil.CupdateWorkflowName(w.Name)))
	defer span.End()

	log := slog.With(slog.String("workflow", w.Name))
	log.DebugContext(ctx, "Running workflow")

	errs := make([]error, len(w.Jobs))

	outputs := NewMap[string, any]()

	done := make([]chan struct{}, len(w.Jobs))
	for i := range w.Jobs {
		done[i] = make(chan struct{})
	}

	// Set all known information for all jobs and steps as we want to report all
	// jobs and steps even if they haven't run (skipped or premature failures)
	// NOTE: workflowRun is generally safe for concurrent writes as accessing a
	// slice by index is thread-safe
	workflowRun := models.WorkflowRun{
		Started: time.Now(),
		Result:  models.WorkflowRunResultSucceeded,
		Jobs:    make([]models.JobRun, len(w.Jobs)),
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		workflowRun.TraceID = spanCtx.TraceID().String()
	}

	for i, job := range w.Jobs {
		workflowRun.Jobs[i] = models.JobRun{
			Result:    models.JobRunResultSkipped,
			Steps:     make([]models.StepRun, len(job.Steps)),
			DependsOn: append([]string{}, job.DependsOn...),
			JobID:     job.ID,
			JobName:   job.Name,
		}

		for j, step := range job.Steps {
			workflowRun.Jobs[i].Steps[j] = models.StepRun{
				Result:   models.StepRunResultSkipped,
				StepName: step.Name,
			}
		}

		// TODO: How do we represent post run steps?
		// Everything has room for it except that we don't really know how to
		// address the jobs so that we can assign to them here?
	}

	var wg sync.WaitGroup
	for i := range w.Jobs {
		job := w.Jobs[i]
		log := log.With(slog.String("job", job.Name))

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(done[i])

			if len(job.DependsOn) > 0 {
				log.DebugContext(ctx, "Waiting for dependencies to complete before starting job")
				for _, dependency := range job.DependsOn {
					index := -1
					for i := range len(w.Jobs) {
						if w.Jobs[i].ID == dependency {
							index = i
							break
						}
					}

					if index != -1 {
						select {
						case <-ctx.Done():
							log.WarnContext(ctx, "Timed out waiting for dependant job to complete", slog.String("dependency", dependency))
							workflowRun.Jobs[i].Result = models.JobRunResultFailed
							errs[i] = ctx.Err()
							return
						case <-done[index]:
							if errs[index] != nil && errs[index] != ErrSkipped {
								log.WarnContext(ctx, "Skipping job as dependent job failed", slog.String("dependency", dependency))
								// Propagate error so that jobs fail if a dependency's
								// dependency fails.
								errs[i] = ErrDependentJobFailed
								return
							}
							// Do nothing
						}
					}
				}
			}

			ctx := Context{
				Context: ctx,

				Workflow:    w,
				WorkflowRun: workflowRun,
				Job:         job,
				JobIndex:    i,

				Outputs: outputs.Clone(),
			}

			started := time.Now()
			ctx, err := job.Run(ctx)
			if err == ErrSkipped {
				errs[i] = ErrSkipped
				return
			}

			workflowRun.Jobs[i].Started = started
			workflowRun.Jobs[i].DurationSeconds = time.Since(workflowRun.Jobs[i].Started).Seconds()

			if err != nil {
				errs[i] = err
				workflowRun.Jobs[i].Result = models.JobRunResultFailed
				return
			}

			workflowRun.Jobs[i].Result = models.JobRunResultSucceeded

			if job.ID != "" {
				for k, v := range ctx.Outputs.Iter() {
					outputs.Set("job."+job.ID+"."+k, v)
				}
			}
		}()
	}

	wg.Wait()

	// Remove unnecessary errors
	errs = slices.DeleteFunc(errs, func(err error) bool {
		return err == ErrDependentJobFailed || err == ErrSkipped
	})

	err := errors.Join(errs...)
	if err == nil {
		span.SetStatus(codes.Ok, "")
	} else {
		span.SetStatus(codes.Error, "One or more jobs failed")
		workflowRun.Result = models.WorkflowRunResultFailed
	}
	return workflowRun, err
}
