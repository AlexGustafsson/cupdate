package gitlab

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

func TestGetProjectContainerRepositories(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetProjectContainerRepositories(context.TODO(), "arm-research/smarter/smarter-device-manager")
	require.NoError(t, err)

	fmt.Printf("%+v\n", res)
}

func TestGetProjectContainerRepositoryTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetProjectContainerRepositoryTags(context.TODO(), "gid://gitlab/ContainerRepository/1080664")
	require.NoError(t, err)

	fmt.Printf("%+v\n", res)
}

func TestClientGetTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}
	ref, err := oci.ParseReference("registry.gitlab.com/arm-research/smarter/smarter-device-manager:v1.20.10")
	require.NoError(t, err)
	tags, err := client.GetTags(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(tags)
}
