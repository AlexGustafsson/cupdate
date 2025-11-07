package static

import (
	"context"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlatformGraph(t *testing.T) {
	staticPlatform := &Platform{
		FilePath: "./testdata/references.txt",
	}

	actual, err := staticPlatform.Graph(context.TODO())
	require.NoError(t, err)

	expected := platform.NewGraph()
	expected.InsertTree(platform.ImageNode{
		Reference: oci.Reference{
			Domain:    "docker.io",
			Path:      "rhasspy/wyoming-whisper",
			HasTag:    true,
			Tag:       "2.5.0",
			HasDigest: true,
			Digest:    "sha256:0d78ad506e450fb113616650b7328233385905e2f2ed07fa59221012144500e3",
		},
	})
	expected.InsertTree(platform.ImageNode{
		Reference: oci.Reference{
			Domain:    "docker.io",
			Path:      "victoriametrics/victoria-metrics",
			HasTag:    true,
			Tag:       "v1.128.0",
			HasDigest: true,
			Digest:    "sha256:c27e736a8aff888cf30c4f20ec648b767358694993d87e89afc6bf80f28991da",
		},
	})
	expected.InsertTree(platform.ImageNode{
		Reference: oci.Reference{
			Domain:    "docker.io",
			Path:      "homeassistant/home-assistant",
			HasTag:    true,
			Tag:       "latest",
			HasDigest: false,
		},
	})
	expected.InsertTree(platform.ImageNode{
		Reference: oci.Reference{
			Domain:    "docker.io",
			Path:      "library/nginx",
			HasTag:    true,
			Tag:       "latest",
			HasDigest: false,
		},
	})
	expected.InsertTree(platform.ImageNode{
		Reference: oci.Reference{
			Domain:    "ghcr.io",
			Path:      "jmbannon/ytdl-sub",
			HasTag:    true,
			Tag:       "latest",
			HasDigest: false,
		},
	})

	assert.Equal(t, expected, actual)
}
