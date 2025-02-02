package workflow

import "errors"

// ErrSkipped is returned when a Job or Step is skipped when invoked.
var ErrSkipped = errors.New("skipped")
