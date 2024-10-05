package workflow

import (
	"fmt"
	"log/slog"
	"strings"
)

type Job struct {
	ID        string
	Name      string
	Steps     []Step
	DependsOn []string
	If        Condition
}

func (j Job) Run(ctx Context) (map[string]any, error) {
	log := slog.With(slog.String("workflow", ctx.Workflow.Name), slog.String("job", ctx.Job.Name))
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

	outputs := make(map[string]any)
	var jobErr error

	log.Debug("Running job steps")
	for i := range j.Steps {
		step := j.Steps[i]
		log := log.With(slog.String("step", step.Name))

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
				return nil, err
			}
		}

		if step.Main != nil && shouldRun {
			log.Debug("Running step")
			command, err := step.Main(ctx)
			if err != nil {
				log.Warn("Job step failed", slog.Any("error", err))
				jobErr = err
				continue
			}

			log.Debug("Step ran successfully")

			// Run side effect
			if command != nil {
				command(ctx)
			}
		}
	}

	log.Debug("Running post steps")
	for _, step := range j.Steps {
		log := log.With(slog.String("step", step.Name))

		shouldRun := jobErr == nil
		if step.PostIf != nil {
			var err error
			shouldRun, err = testCondition(ctx, step.PostIf)
			if err != nil {
				return nil, err
			}
		}

		if step.Post != nil && shouldRun {
			log.Debug("Running step")
			if err := step.Post(ctx); err != nil {
				log.Warn("Job post step failed", slog.Any("error", err))
				jobErr = err
				continue
			}

			log.Debug("Post step ran successfully")
		}
	}

	if jobErr != nil {
		return nil, fmt.Errorf("job failed due to one or more errors")
	}

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
