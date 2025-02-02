package workflow

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
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

	log := slog.With(slog.String("workflow", w.Name))
	log.DebugContext(ctx, "Running workflow")

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
		log := log.With(slog.String("job", job.Name))

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(done[i])

			if len(job.DependsOn) > 0 {
				log.DebugContext(ctx, "Waiting for dependencies to complete before starting job")
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
								log.WarnContext(ctx, "Skipping job as dependent job failed", slog.String("dependency", dependency))
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

			ctx, err := job.Run(ctx)
			if err == ErrSkipped {
				return
			} else if err != nil {
				errs[i] = err
				return
			}

			mutex.Lock()
			if job.ID != "" {
				for k, v := range ctx.Outputs {
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
