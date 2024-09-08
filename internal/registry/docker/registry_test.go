package docker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/k8s-image-feed/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryGet(t *testing.T) {
	expected := &registry.Image{
		Name:         "postgres",
		Version:      "13.15-bullseye",
		Published:    time.Date(2024, 05, 15, 17, 8, 59, 365138000, time.UTC),
		Digest:       "sha256:d9b56dd5a190a14dc9368383eaa6b846168442c18a97e15f814727273f82d9ce",
		ReleaseNotes: "",
	}

	var registry Registry
	actual, err := registry.Get(context.TODO(), "postgres", "sha256:d9b56dd5a190a14dc9368383eaa6b846168442c18a97e15f814727273f82d9ce")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
func TestRegistryGetManifest(t *testing.T) {
	expected := &Manifest{}

	var registry Registry
	actual, err := registry.GetManifests(context.TODO(), "postgres", "12-alpine")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestRegistryGetLatestVersion(t *testing.T) {
	var registry Registry
	image, err := registry.GetLatestVersion(context.TODO(), "renovate/renovate", "38.70.2")
	require.NoError(t, err)

	fmt.Println(image)
}

func TestRegistryGetRepository(t *testing.T) {
	var registry Registry
	repository, err := registry.GetRepository(context.TODO(), "homeassistant", "home-assistant")
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)
}
