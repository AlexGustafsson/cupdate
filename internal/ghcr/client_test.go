package ghcr

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

func TestOCI(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	ociClient := &oci.Client{
		Client:   client.Client,
		AuthFunc: client.HandleAuth,
	}

	references := []string{
		"ghcr.io/jmbannon/ytdl-sub:2024.10.09",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := oci.ParseReference(reference)
			require.NoError(t, err)

			manifest, err := ociClient.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)

			// Rewrite ref to pin digest of manifst
			ref.HasTag = false
			ref.Tag = ""
			ref.HasDigest = true
			switch m := manifest.(type) {
			case *oci.ImageManifest:
				ref.Digest = m.Digest
			case *oci.ImageIndex:
				ref.Digest = m.Digest
			}

			// Expect it to exist
			manifest, err = ociClient.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)
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

	ref, err := oci.ParseReference("ghcr.io/jmbannon/ytdl-sub")
	require.NoError(t, err)

	ociClient := &oci.Client{
		Client:   client.Client,
		AuthFunc: client.HandleAuth,
	}

	manifests, err := ociClient.GetAnnotations(context.TODO(), ref, nil)
	require.NoError(t, err)

	fmt.Println(manifests, manifests == nil)
}
