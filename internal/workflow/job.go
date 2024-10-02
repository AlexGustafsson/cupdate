package workflow

import "log/slog"

type Job struct {
	ID        string
	Name      string
	Steps     []Step
	DependsOn []string
	ShouldRun func(ctx Context) (bool, error)
}

func (j Job) Run(ctx Context) (map[string]any, error) {
	log := slog.With(slog.String("workflow", ctx.Workflow), slog.String("job", ctx.Job))

	if j.ShouldRun != nil {
		shouldRun, err := j.ShouldRun(ctx)
		if err != nil {
			log.Error("Failed to identify if job should run", slog.Any("error", err))
			return nil, err
		}

		if !shouldRun {
			log.Debug("Skipping job")
			return nil, nil
		}
	}

	outputs := make(map[string]any)

	log.Debug("Running job")
	for i := range j.Steps {
		step := j.Steps[i]

		log := log.With(slog.String("step", step.Name()))
		log.Debug("Running step")
		ctx := Context{
			Context: ctx.Context,

			Workflow: ctx.Workflow,
			Job:      ctx.Job,
			Step:     step.Name(),

			Outputs: outputs,
			Inputs:  make(map[string]string),
		}
		stepOutputs, err := step.Run(ctx)
		if err != nil {
			log.Warn("Job step failed", slog.Any("error", err))
			return nil, err
		}
		log.Debug("Step ran successfully")

		if step.ID() != "" {
			for k, v := range stepOutputs {
				outputs["step."+step.ID()+"."+k] = v
			}
		}
	}

	return outputs, nil
}
