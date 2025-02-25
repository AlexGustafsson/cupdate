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

func TestGetScan(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ref, err := oci.ParseReference("quay.io/openshift-release-dev/ocp-release@sha256:7708f832ae02919f2cdb2798fdbc64e17ce7a576d1e3baabdd78a000d2d62f40")
	require.NoError(t, err)

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	scan, err := client.GetScan(context.TODO(), ref)
	require.NoError(t, err)

	fmt.Printf("%+v\n", scan)
}
