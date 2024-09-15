package github

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientGetRelease(t *testing.T) {
	c := &Client{}

	release, err := c.GetRelease(context.TODO(), "renovatebot", "renovate", "38.80.0")
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}
