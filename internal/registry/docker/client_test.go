package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGet(t *testing.T) {
	expected := &registry.Image{
		Name:         "postgres",
		Version:      "13.15-bullseye",
		Published:    time.Date(2024, 05, 15, 17, 8, 59, 365138000, time.UTC),
		Digest:       "sha256:d9b56dd5a190a14dc9368383eaa6b846168442c18a97e15f814727273f82d9ce",
		ReleaseNotes: "",
	}

	var client Client
	actual, err := client.Get(context.TODO(), "postgres", "sha256:d9b56dd5a190a14dc9368383eaa6b846168442c18a97e15f814727273f82d9ce")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
func TestClientGetManifest(t *testing.T) {
	expected := &oci.Manifest{}

	var client Client
	ref, err := reference.ParseNormalizedNamed("postgres:12-alpine")
	require.NoError(t, err)
	actual, err := client.GetManifests(context.TODO(), ref.(reference.NamedTagged))
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestClientGetLatestVersion(t *testing.T) {
	var client Client
	image, err := client.GetLatestVersion(context.TODO(), "renovate/renovate", "38.70.2")
	require.NoError(t, err)

	fmt.Println(image)
}

func TestClientGetRepository(t *testing.T) {
	var client Client
	repository, err := client.GetRepository(context.TODO(), "mongo")
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)

	json.NewEncoder(os.Stdout).Encode(repository)
}
