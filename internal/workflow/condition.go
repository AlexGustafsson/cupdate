package workflow

import (
	"fmt"
)

// Condition controls whether or not a job or step (or post step) should run.
// A condition can be either a [ConditionFunc], or a string which references a
// value retrievable using [GetValue].
type Condition = any

// ConditionFunc is an adapter to allow the use of ordinary functions as
// Conditions.
type ConditionFunc = func(ctx Context) (bool, error)

// Always will always run.
func Always(ctx Context) (bool, error) {
	return true, nil
}

// ValueExists will evaluate to true if the given key has a value.
func ValueExists(key string) Condition {
	return ConditionFunc(func(ctx Context) (bool, error) {
		_, ok := GetAnyValue(ctx, key)
		return ok, nil
	})
}

// testCondition tests whether or not a condition is fulfilled.
func testCondition(ctx Context, condition Condition) (bool, error) {
	conditionFunc, ok := condition.(ConditionFunc)
	if ok {
		return conditionFunc(ctx)
	}

	switch cond := condition.(type) {
	case string:
		// The condition was a reference to a boolean
		return GetValue[bool](ctx, cond)
	default:
		return false, fmt.Errorf("invalid condition type %T expected string or %T", condition, conditionFunc)
	}
}
