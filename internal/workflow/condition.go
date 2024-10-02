package workflow

import "fmt"

// Condition controls whether or not a job or step (or post step) should run.
// A condition can be either a [ConditionFunc], or a string which references a
// value retrievable using [GetValue].
type Condition = any

// Conditio
type ConditionFunc func(ctx Context) (bool, error)

// Always will always run.
func Always(ctx Context) (bool, error) {
	return true, nil
}

// testCondition tests whether or not a condition is fulfilled.
func testCondition(ctx Context, condition Condition) (bool, error) {
	conditionFunc, ok := condition.(ConditionFunc)
	if ok {
		return conditionFunc(ctx)
	}

	switch cond := condition.(type) {
	case string:
		// The condition was a reference to a
		return GetValue[bool](ctx, cond)
	default:
		return false, fmt.Errorf("invalid condition")
	}
}
