package workflow

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type Workflow struct {
	Name string
	Jobs []Job
}

func (w Workflow) Run(ctx context.Context) error {
	log := slog.With(slog.String("workflow", w.Name))
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
		log := log.With(slog.String("job", job.Name))

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(done[i])

			if len(job.DependsOn) > 0 {
				log.Debug("Waiting for dependencies to complete before starting job")
				for _, dependant := range job.DependsOn {
					index := -1
					for i := 0; i < len(w.Jobs); i++ {
						if w.Jobs[i].ID == dependant {
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
					outputs["jobs."+job.ID+"."+k] = v
				}
			}
			mutex.Unlock()
		}()
	}

	wg.Wait()
	return errors.Join(errs...)
}

// TODO: Build a "dot" graph for debugging / visualization?
// func (w Workflow) Describe() string {
//
// }
