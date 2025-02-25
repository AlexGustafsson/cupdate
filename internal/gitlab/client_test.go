package gitlab

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
		"registry.gitlab.com/arm-research/smarter/smarter-device-manager",
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

func TestGetRepositoryDescription(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetRepositoryDescription(context.TODO(), "arm-research/smarter/smarter-device-manager")
	require.NoError(t, err)

	fmt.Printf("%+v\n", res)
}

func TestGetRepositoryREADMEBlob(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	res, err := client.GetRepositoryREADMEBlob(context.TODO(), "/arm-research/smarter/smarter-device-manager/-/blob/master/README.md")
	require.NoError(t, err)

	fmt.Printf("%+s\n", res.Raw)
}
