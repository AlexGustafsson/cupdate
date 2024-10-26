package quay

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

func TestClientGetLatestVersion(t *testing.T) {
	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}
	ref, err := oci.ParseReference("quay.io/jetstack/cert-manager-acmesolver:v1.16.0")
	require.NoError(t, err)
	actual, err := client.GetLatestVersion(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(actual)
}
