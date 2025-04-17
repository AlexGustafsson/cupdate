package workflow

import (
	"log/slog"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Step represents a single unit of an operation.
// The implementation mimics GitHub actions / workflows.
type Step struct {
	ID   string
	Name string

	// Inputs holds a map of named inputs that can be used by the step.
	Inputs map[string]Input

	Main func(ctx Context) (Command, error)
	If   Condition

	Post   func(ctx Context) error
	PostIf Condition
}

func (s Step) WithID(id string) Step {
	s.ID = id
	return s
}

func (s Step) WithCondition(condition Condition) Step {
	if s.If != nil {
		panic("step already has a condition")
	}
	s.If = condition
	return s
}

func (s Step) With(key string, input Input) Step {
	if s.Inputs == nil {
		s.Inputs = make(map[string]any)
	}
	s.Inputs[key] = input
	return s
}

// Run runs a step.
// Returns the modified context.
// Returns [ErrSkipped] if the job was not run.
// It's the caller's responsibility to ensure [Step.Main] is defined.
func (s Step) Run(ctx Context) (Context, error) {
	log := slog.With(slog.String("workflow", ctx.Workflow.Name), slog.String("job", ctx.Job.Name), slog.String("step", ctx.Step.Name))

	shouldRun := ctx.Error == nil
	if s.If != nil {
		var err error
		shouldRun, err = testCondition(ctx, s.If)
		if err != nil {
			return ctx, err
		}
	}

	if !shouldRun {
		log.DebugContext(ctx, "Skipped step")
		return ctx, ErrSkipped
	}

	var stepSpan trace.Span
	ctx.Context, stepSpan = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, otelutil.CupdateWorkflowStepRunSpanName, trace.WithAttributes(otelutil.CupdateWorkflowStepName(s.Name)))
	defer stepSpan.End()

	log.DebugContext(ctx, "Running step")

	command, err := s.Main(ctx)
	if err != nil {
		log.WarnContext(ctx, "Job step failed", slog.Any("error", err))
		stepSpan.SetStatus(codes.Error, "Step failed")
		stepSpan.End()
		return ctx, err
	}

	log.DebugContext(ctx, "Step ran successfully")

	// Run side effect
	if command != nil {
		command(ctx)
	}

	stepSpan.SetStatus(codes.Ok, "")

	return ctx, nil
}

// RunPost runs a step's post step.
// Returns [ErrSkipped] if the job was not run.
// It's the caller's responsibility to ensure [Step.Post] is defined.
func (s Step) RunPost(ctx Context) error {
	log := slog.With(slog.String("workflow", ctx.Workflow.Name), slog.String("job", ctx.Job.Name), slog.String("step", ctx.Step.Name))

	shouldRun := ctx.Error == nil
	if s.PostIf != nil {
		var err error
		shouldRun, err = testCondition(ctx, s.PostIf)
		if err != nil {
			return err
		}
	}

	if !shouldRun {
		log.DebugContext(ctx, "Skipped post step")
		return ErrSkipped
	}

	var stepSpan trace.Span
	ctx.Context, stepSpan = otel.Tracer(otelutil.DefaultScope).Start(ctx.Context, otelutil.CupdateWorkflowStepPostRunSpanName, trace.WithAttributes(otelutil.CupdateWorkflowStepName(s.Name)))
	defer stepSpan.End()

	log.DebugContext(ctx, "Running post step")

	err := s.Post(ctx)
	if err != nil {
		log.WarnContext(ctx, "Job post step failed", slog.Any("error", err))
		stepSpan.SetStatus(codes.Error, "Post step failed")
		stepSpan.End()
		return err
	}

	log.DebugContext(ctx, "Post step ran successfully")

	stepSpan.SetStatus(codes.Ok, "")

	return nil
}

type StepFunc func(ctx Context) (Command, error)

func Run(f StepFunc) Step {
	return Step{
		Main: f,
	}
}
