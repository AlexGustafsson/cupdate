package ghcr

import (
	"context"
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
	ref, err := oci.ParseReference("ghcr.io/jmbannon/ytdl-sub:2024.10.09")
	require.NoError(t, err)
	actual, err := client.GetManifests(context.TODO(), ref)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
