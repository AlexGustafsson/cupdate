package github

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/stretchr/testify/require"
)

func TestClientGetRelease(t *testing.T) {
	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	release, err := c.GetRelease(context.TODO(), "renovatebot", "renovate", "38.80.0")
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}

func TestClientGetDescription(t *testing.T) {
	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	release, err := c.GetDescription(context.TODO(), "renovatebot", "renovate")
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}

func TestClientGetPackage(t *testing.T) {
	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ref, err := oci.ParseReference("ghcr.io/jmbannon/ytdl-sub")
	require.NoError(t, err)

	release, err := c.GetPackage(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}
