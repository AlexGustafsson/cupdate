package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetManifest(t *testing.T) {
	expected := &oci.Manifest{}

	client := &Client{
		Client: httputil.NewClient(cache.NewInMemoryCache(), 24*time.Hour),
	}
	ref, err := oci.ParseReference("postgres:12-alpine")
	require.NoError(t, err)
	actual, err := client.GetManifests(context.TODO(), ref)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestClientGetLatestVersion(t *testing.T) {
	client := &Client{
		Client: httputil.NewClient(cache.NewInMemoryCache(), 24*time.Hour),
	}
	ref, err := oci.ParseReference("renovate/renovate:38.70.2")
	require.NoError(t, err)
	image, err := client.GetLatestVersion(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(image)
}

func TestClientGetRepository(t *testing.T) {
	client := &Client{
		Client: httputil.NewClient(cache.NewInMemoryCache(), 24*time.Hour),
	}
	ref, err := oci.ParseReference("mongo")
	require.NoError(t, err)
	repository, err := client.GetRepository(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)

	json.NewEncoder(os.Stdout).Encode(repository)
}
