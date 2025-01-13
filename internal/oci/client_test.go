package oci

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/require"
)

func TestClientGetManifest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ref, err := ParseReference("k8s.gcr.io/pause")
	require.NoError(t, err)

	actual, err := client.GetManifests(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Printf("%+v\n", actual)
}
