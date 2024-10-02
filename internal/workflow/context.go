package workflow

import (
	"context"
	"fmt"
)

var _ context.Context = (*Context)(nil)

type Context struct {
	context.Context

	Workflow string
	Job      string
	Step     string

	Outputs map[string]any
	Inputs  map[string]string
}

func (c Context) Output(key string) (any, bool) {
	v, ok := c.Outputs[key]
	return v, ok
}

func (c Context) Input(key string) (any, error) {
	outputKey, ok := c.Inputs[key]
	if !ok {
		return nil, fmt.Errorf("no value specified for input: %s", key)
	}

	output, ok := c.Output(outputKey)
	if !ok {
		return nil, fmt.Errorf("no value stored for input: %s", outputKey)
	}

	return output, nil
}
