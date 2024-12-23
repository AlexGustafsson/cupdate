package github

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/require"
)

func TestClientGetRelease(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	release, err := c.GetRelease(context.TODO(), "renovatebot", "renovate", "38.80.0")
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}

func TestClientGetDescription(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	release, err := c.GetDescription(context.TODO(), "renovatebot", "renovate")
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}

func TestClientGetPackage(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ref, err := oci.ParseReference("ghcr.io/alexgustafsson/srdl")
	require.NoError(t, err)

	release, err := c.GetPackage(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Printf("%+v\n", release)
}

func TestClientGetREADME(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	c := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	readme, err := c.GetREADME(context.TODO(), "https://github.com/users/AlexGustafsson/packages/container/srdl/307129679/readme?is_package_page=true")
	require.NoError(t, err)

	fmt.Printf("%s\n", readme)
}
