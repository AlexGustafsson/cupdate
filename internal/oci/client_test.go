package oci

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientGetManifest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"k8s.gcr.io/pause",
		"quay.io/jetstack/cert-manager-startupapicheck:v1.16.2",
		"registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.14.0",
		"gcr.io/zenika-hub/alpine-chrome:123",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			manifest, err := client.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)

			// Rewrite ref to pin digest of manifst
			ref.HasTag = false
			ref.Tag = ""
			ref.HasDigest = true
			switch m := manifest.(type) {
			case *ImageManifest:
				ref.Digest = m.Digest
			case *ImageIndex:
				ref.Digest = m.Digest
			}

			// Expect it to exist
			manifest, err = client.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)
		})
	}
}

func TestClientHeadBlob(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	testCases := []struct {
		Reference string
		Digest    string
	}{
		{
			Reference: "k8s.gcr.io/pause",
			Digest:    "sha256:350b164e7ae1dcddeffadd65c76226c9b6dc5553f5179153fb0e36b78f2a5e06",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Reference, func(t *testing.T) {
			ref, err := ParseReference(testCase.Reference)
			require.NoError(t, err)

			info, err := client.HeadBlob(context.TODO(), ref, testCase.Digest)
			require.NoError(t, err)
			fmt.Printf("%+v\n", info)
		})
	}
}

func TestClientGetAnnotations(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"quay.io/jetstack/cert-manager-startupapicheck:v1.16.2",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			annotations, err := client.GetAnnotations(context.TODO(), ref, nil)
			require.NoError(t, err)
			assert.NotNil(t, annotations)
			fmt.Printf("%+v\n", annotations)
		})
	}
}

func TestClientGetTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"k8s.gcr.io/pause",
		"quay.io/jetstack/cert-manager-startupapicheck",
		"registry.k8s.io/kube-state-metrics/kube-state-metrics",
		"gcr.io/zenika-hub/alpine-chrome",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			tags, err := client.GetTags(context.TODO(), ref, nil)
			require.NoError(t, err)
			assert.NotNil(t, tags)
			fmt.Printf("%+v\n", tags)
		})
	}
}
