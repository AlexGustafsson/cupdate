package workflow

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"github.com/AlexGustafsson/cupdate/internal/slogutil"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Workflow struct {
	Name string
	Jobs []Job
}

func (w Workflow) Run(ctx context.Context) error {
	ctx, span := otel.Tracer(otelutil.DefaultScope).Start(ctx, otelutil.CupdateWorkflowRunSpanName, trace.WithAttributes(otelutil.CupdateWorkflowName(w.Name)))
	defer span.End()

	log := slog.With(slog.String("workflow", w.Name)).With(slogutil.Context(ctx))
	log.Debug("Running workflow")

	var mutex sync.Mutex
	errs := make([]error, len(w.Jobs))
	outputs := make(map[string]any)

	done := make([]chan struct{}, len(w.Jobs))
	for i := range w.Jobs {
		done[i] = make(chan struct{})
	}

	var wg sync.WaitGroup
	for i := range w.Jobs {
		job := w.Jobs[i]
		log := log.With(slog.String("job", job.Name)).With(slogutil.Context(ctx))

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(done[i])

			if len(job.DependsOn) > 0 {
				log.Debug("Waiting for dependencies to complete before starting job")
				for _, dependency := range job.DependsOn {
					index := -1
					for i := 0; i < len(w.Jobs); i++ {
						if w.Jobs[i].ID == dependency {
							index = i
							break
						}
					}

					if index != -1 {
						select {
						case <-ctx.Done():
							mutex.Lock()
							errs[i] = ctx.Err()
							mutex.Unlock()
						case <-done[index]:
							if errs[index] != nil {
								log.Warn("Skipping job as dependent job failed", slog.String("dependency", dependency))
								errs[i] = fmt.Errorf("failed to run job - dependent job failed: %w", errs[index])
								return
							}
							// Do nothing
						}
					}
				}
			}

			ctx := Context{
				Context: ctx,

				Workflow: w,
				Job:      job,

				Outputs: outputs,
			}

			jobOutputs, err := job.Run(ctx)

			mutex.Lock()
			errs[i] = err
			if job.ID != "" {
				for k, v := range jobOutputs {
					outputs["job."+job.ID+"."+k] = v
				}
			}
			mutex.Unlock()
		}()
	}

	wg.Wait()

	err := errors.Join(errs...)
	if err == nil {
		span.SetStatus(codes.Ok, "")
	} else {
		span.SetStatus(codes.Error, "One or more jobs failed")
	}
	return err
}

func (w Workflow) Describe() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, `---
title: %s
---
flowchart LR
`, w.Name)

	fmt.Fprintf(&builder, "start[Start] --> job.0\n")
	fmt.Fprintf(&builder, "job.%d --> stop[Stop]\n", len(w.Jobs)-1)

	for i, job := range w.Jobs {
		builder.WriteString(job.Describe(fmt.Sprintf("job.%d", i)))

		for _, dependency := range job.DependsOn {
			for j, job := range w.Jobs {
				if job.ID == dependency {
					fmt.Fprintf(&builder, "job.%d -- depends on --> job.%d\n", i, j)
					break
				}
			}
		}
	}

	return builder.String()
}
