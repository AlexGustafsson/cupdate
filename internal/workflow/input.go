package workflow

import (
	"fmt"
)

// Input is an input value to a step.
// A value can be either a [Ref], or any verbatim value.
type Input = any

// GetInput returns a value or output in the ctx.
// If a value does not exist, the type's zero value is returned, with a nil
// error unless required is true.
func GetInput[T any](ctx Context, name string, required bool) (T, error) {
	var ret T
	v, err := GetAnyInput(ctx, name, required)
	if err != nil {
		return ret, err
	}

	var ok bool
	ret, ok = v.(T)
	if !ok {
		return ret, fmt.Errorf("invalid type %T for input %s of type %T", v, name, ret)
	}

	return ret, nil
}

// GetAnyInput returns a value or output in the ctx.
// If a value does not exist, nil is returned, with a nil error unless required
// is true.
func GetAnyInput(ctx Context, name string, required bool) (any, error) {
	if ctx.Step.Inputs == nil {
		if required {
			return nil, fmt.Errorf("missing required input: %s", name)
		}
		return nil, nil
	}

	v, ok := ctx.Step.Inputs[name]
	if !ok {
		if required {
			return nil, fmt.Errorf("missing required input: %s", name)
		}
		return nil, nil
	}

	switch v := v.(type) {
	case Ref:
		ret, ok := GetAnyValue(ctx, v.Key)
		if !ok {
			if required {
				return nil, fmt.Errorf("missing required input: %s", name)
			}
			return nil, nil
		}

		return ret, nil
	case Func:
		return v.Func(ctx)
	default:
		return v, nil
	}
}

// Ref is an [Input] that refers to a value retrievable by [GetValue].
type Ref struct {
	Key string
}

// Func is an [Input] that is invoked to get a value.
type Func struct {
	Func func(ctx Context) (any, error)
}
