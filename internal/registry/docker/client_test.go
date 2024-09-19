package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/distribution/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	ref, err := reference.ParseNormalizedNamed("renovate/renovate:38.70.2")
	require.NoError(t, err)
	image, err := client.GetLatestVersion(context.TODO(), ref.(reference.NamedTagged))
	require.NoError(t, err)

	fmt.Println(image)
}

func TestClientGetRepository(t *testing.T) {
	var client Client
	ref, err := reference.ParseNormalizedNamed("mongo")
	require.NoError(t, err)
	repository, err := client.GetRepository(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(repository.FullDescription)

	json.NewEncoder(os.Stdout).Encode(repository)
}
