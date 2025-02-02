package workflow

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Job struct {
	ID        string
	Name      string
	Steps     []Step
	DependsOn []string
	If        Condition
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

	outputs := maps.Clone(ctx.Outputs)

	errs := make([]error, len(j.Steps)*2)

	log.DebugContext(ctx, "Running job steps")
	for i, step := range j.Steps {
		if step.Main == nil {
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

		ctx, err := step.Run(ctx)
		if err == ErrSkipped {
			continue
		} else if err != nil {
			errs[i] = err
			continue
		}

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

func (j Job) Describe(namespace string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "subgraph %s [%s]\n", namespace, j.Name)

	for i, step := range j.Steps {
		builder.WriteString(step.Describe(fmt.Sprintf("%s.step.%d", namespace, i)))
	}

	for i := 1; i < len(j.Steps); i++ {
		fmt.Fprintf(&builder, "%s.step.%d --> %s.step.%d\n", namespace, i-1, namespace, i)
	}

	fmt.Fprintf(&builder, "end\n")

	return builder.String()
}
