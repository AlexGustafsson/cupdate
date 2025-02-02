package workflow

import "fmt"

// Command represents a side effect of a step.
// A step can return a command, which will then be executed after the stes
// completes. This enables a step to define outputs, values and more.
type Command func(ctx Context)

func Batch(commands ...Command) Command {
	return func(ctx Context) {
		for _, command := range commands {
			command(ctx)
		}
	}
}

// SetOutput returns a [Command] that will store an output.
// Setting an output from a step without an ID is a no-op.
// Outputs are named depending on their scope. Within the same job, a step will
// be named like so:
//
//	"step.<stepId>.<key>"
//
// From another job, the output can be referenced like so:
//
//	"job.<jobId>.step.<stepId>.<key>"
func SetOutput(key string, value any) Command {
	return func(ctx Context) {
		if ctx.Step.ID == "" {
			return
		}

		k := fmt.Sprintf("step.%s.%s", ctx.Step.ID, key)
		ctx.Outputs[k] = value
	}
}
