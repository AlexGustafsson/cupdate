package quay

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

func TestClientGetTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}
	ref, err := oci.ParseReference("quay.io/jetstack/cert-manager-acmesolver:v1.16.0")
	require.NoError(t, err)
	tags, err := client.GetTags(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Println(tags)
}
