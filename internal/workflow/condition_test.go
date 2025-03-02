package workflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlways(t *testing.T) {
	shouldRun, err := Always(Context{})
	assert.True(t, shouldRun)
	assert.NoError(t, err)
}

func TestValueExists(t *testing.T) {
	ctx := Context{
		Outputs: map[string]any{
			"foo": "bar",
		},
	}

	condition := ValueExists("foo")
	shouldRun, err := condition(ctx)
	assert.True(t, shouldRun)
	assert.NoError(t, err)

	condition = ValueExists("bar")
	shouldRun, err = condition(ctx)
	assert.False(t, shouldRun)
	assert.NoError(t, err)

	// Job has already failed
	ctx.Error = fmt.Errorf("some error")
	condition = ValueExists("foo")
	shouldRun, err = condition(ctx)
	assert.False(t, shouldRun)
	assert.NoError(t, err)
}
