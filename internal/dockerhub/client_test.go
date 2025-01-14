package dockerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetManifest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	expected := &oci.Manifest{}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ref, err := oci.ParseReference("postgres:12-alpine")
	require.NoError(t, err)

	ociClient := &oci.Client{
		Client:   client.Client,
		AuthFunc: client.HandleAuth,
	}

	actual, err := ociClient.GetManifests(context.TODO(), ref)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestClientGetAnnotations(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ref, err := oci.ParseReference("homeassistant/home-assistant")
	require.NoError(t, err)

	ociClient := &oci.Client{
		Client:   client.Client,
		AuthFunc: client.HandleAuth,
	}

	manifests, err := ociClient.GetAnnotations(context.TODO(), ref, nil)
	require.NoError(t, err)

	fmt.Println(manifests, manifests == nil)
}

func TestClientGetRepository(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}
	ref, err := oci.ParseReference("mongo")
	require.NoError(t, err)
	repository, err := client.GetRepository(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)

	json.NewEncoder(os.Stdout).Encode(repository)
}

func TestGetVulnerabilityReport(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	report, err := client.GetVulnerabilityReport(context.TODO(), "traefik", "sha256:ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078")
	require.NoError(t, err)

	json.NewEncoder(os.Stdout).Encode(report)
}

func TestGetTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ociClient := oci.Client{
		Client:   client.Client,
		AuthFunc: client.HandleAuth,
	}

	ref, err := oci.ParseReference("mongo")
	require.NoError(t, err)

	tags, err := ociClient.GetTags(context.TODO(), ref, &oci.GetTagsOptions{
		Count:    300,
		AllPages: true,
	})
	require.NoError(t, err)

	fmt.Println(tags)
}
