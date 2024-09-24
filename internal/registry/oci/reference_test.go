package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReference(t *testing.T) {
	testCase := []struct {
		Reference string

		ExpectedName   string
		ExpectedString string
		ExpectedDomain string
		ExpectedPath   string
		ExpectedTag    string
		ExpectedDigest string
	}{
		{
			Reference: "mongo",

			ExpectedName:   "mongo",
			ExpectedString: "mongo",
			ExpectedDomain: "docker.io",
			ExpectedPath:   "library/mongo",
		},
		{
			Reference: "mongo:4",

			ExpectedName:   "mongo",
			ExpectedString: "mongo:4",
			ExpectedDomain: "docker.io",
			ExpectedPath:   "library/mongo",
			ExpectedTag:    "4",
		},
		{
			Reference: "library/mongo:4",

			ExpectedName:   "mongo",
			ExpectedString: "mongo:4",
			ExpectedDomain: "docker.io",
			ExpectedPath:   "library/mongo",
			ExpectedTag:    "4",
		},
		{
			Reference: "docker.io/library/mongo:4",

			ExpectedName:   "mongo",
			ExpectedString: "mongo:4",
			ExpectedDomain: "docker.io",
			ExpectedPath:   "library/mongo",
			ExpectedTag:    "4",
		},
		{
			Reference: "ghcr.io/mongo/mongo",

			ExpectedName:   "ghcr.io/mongo/mongo",
			ExpectedString: "ghcr.io/mongo/mongo",
			ExpectedDomain: "ghcr.io",
			ExpectedPath:   "mongo/mongo",
		},
		{
			Reference: "ghcr.io/mongo/mongo:4",

			ExpectedName:   "ghcr.io/mongo/mongo",
			ExpectedString: "ghcr.io/mongo/mongo:4",
			ExpectedDomain: "ghcr.io",
			ExpectedPath:   "mongo/mongo",
			ExpectedTag:    "4",
		},
		{
			Reference: "mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",

			ExpectedName:   "mongo",
			ExpectedString: "mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			ExpectedDomain: "docker.io",
			ExpectedPath:   "library/mongo",
			ExpectedDigest: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			Reference: "ghcr.io/mongo/mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",

			ExpectedName:   "ghcr.io/mongo/mongo",
			ExpectedString: "ghcr.io/mongo/mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			ExpectedDomain: "ghcr.io",
			ExpectedPath:   "mongo/mongo",
			ExpectedDigest: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, testCase := range testCase {
		t.Run(testCase.Reference, func(t *testing.T) {
			r, err := ParseReference(testCase.Reference)
			require.NoError(t, err)

			assert.Equal(t, testCase.ExpectedName, r.Name())
			assert.Equal(t, testCase.ExpectedString, r.String())
			assert.Equal(t, testCase.ExpectedDomain, r.Domain)
			assert.Equal(t, testCase.ExpectedPath, r.Path)

			if testCase.ExpectedTag == "" {
				assert.False(t, r.HasTag)
				assert.Equal(t, "", r.Tag)
			} else {
				assert.True(t, r.HasTag)
				assert.Equal(t, testCase.ExpectedTag, r.Tag)
			}

			if testCase.ExpectedDigest == "" {
				assert.False(t, r.HasDigest)
				assert.Equal(t, "", r.Digest)
			} else {
				assert.True(t, r.HasDigest)
				assert.Equal(t, testCase.ExpectedDigest, r.Digest)
			}
		})
	}
}
