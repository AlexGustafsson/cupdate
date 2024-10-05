package workflow

import (
	"fmt"
	"strings"
)

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

func (s Step) Describe(namespace string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%s[%s]\n", namespace, s.Name)

	return builder.String()
}
