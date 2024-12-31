package workflow

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"github.com/AlexGustafsson/cupdate/internal/slogutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

func (j Job) Run(ctx Context) (map[string]any, error) {
	log := slog.With(slog.String("workflow", ctx.Workflow.Name), slog.String("job", ctx.Job.Name)).With(slogutil.Context(ctx))
	log.Debug("Running job")

	if j.If != nil {
		shouldRun, err := testCondition(ctx, j.If)
		if err != nil {
			log.Error("Failed to identify if job should run", slog.Any("error", err))
			return nil, err
		}

		if !shouldRun {
			log.Debug("Skipping job in accordance to specified condition")
			return nil, nil
		}
	}

	var jobSpan trace.Span
	ctx.Context, jobSpan = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, "cupdate.workflow.job.run")
	jobSpan.SetAttributes(attribute.String("cupdate.workflow.job.name", j.Name))
	defer jobSpan.End()

	outputs := make(map[string]any)
	for k, v := range ctx.Outputs {
		outputs[k] = v
	}

	var jobErr error

	log.Debug("Running job steps")
	for i := range j.Steps {
		step := j.Steps[i]
		log := log.With(slog.String("step", step.Name)).With(slogutil.Context(ctx))

		ctx := Context{
			Context: ctx.Context,

			Workflow: ctx.Workflow,
			Job:      ctx.Job,
			Step:     step,

			Outputs: outputs,

			Error: jobErr,
		}

		shouldRun := jobErr == nil
		if step.If != nil {
			var err error
			shouldRun, err = testCondition(ctx, step.If)
			if err != nil {
				jobSpan.SetStatus(codes.Error, "Condition test failure")
				return nil, err
			}
		}

		if step.Main != nil && shouldRun {
			var stepSpan trace.Span
			ctx.Context, stepSpan = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, "cupdate.workflow.step.run")
			stepSpan.SetAttributes(attribute.String("cupdate.workflow.step.name", step.Name))

			log := log.With(slogutil.Context(ctx))

			log.Debug("Running step")

			command, err := step.Main(ctx)
			if err != nil {
				log.Warn("Job step failed", slog.Any("error", err))
				jobErr = err
				stepSpan.SetStatus(codes.Error, "Step failed")
				stepSpan.End()
				continue
			}

			log.Debug("Step ran successfully")

			// Run side effect
			if command != nil {
				command(ctx)
			}

			stepSpan.SetStatus(codes.Ok, "")
			stepSpan.End()
		}
	}

	log.Debug("Running post steps")
	for i := range j.Steps {
		step := j.Steps[i]
		if step.Post == nil {
			continue
		}

		ctx := Context{
			Context: ctx.Context,

			Workflow: ctx.Workflow,
			Job:      ctx.Job,
			Step:     step,

			Outputs: outputs,

			Error: jobErr,
		}

		var postStepRun trace.Span
		ctx.Context, postStepRun = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, "cupdate.workflow.step.post-run")
		postStepRun.SetAttributes(attribute.String("cupdate.workflow.step.name", step.Name))

		log := log.With(slog.String("step", step.Name)).With(slogutil.Context(ctx))

		shouldRun := jobErr == nil
		if step.PostIf != nil {
			var err error
			shouldRun, err = testCondition(ctx, step.PostIf)
			if err != nil {
				postStepRun.SetStatus(codes.Error, "Condition test failure")
				postStepRun.End()
				return nil, err
			}
		}

		if shouldRun {
			log.Debug("Running post step")
			if err := step.Post(ctx); err != nil {
				log.Warn("Job post step failed", slog.Any("error", err))
				jobErr = err
				postStepRun.SetStatus(codes.Error, "Post step failed")
				postStepRun.End()
				continue
			}

			log.Debug("Post step ran successfully")
			postStepRun.SetStatus(codes.Ok, "")
			postStepRun.End()
		}
	}

	if jobErr != nil {
		jobSpan.SetStatus(codes.Error, "One or more steps failed")
		return nil, fmt.Errorf("job failed due to one or more errors")
	}

	jobSpan.SetStatus(codes.Ok, "")
	return outputs, nil
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
