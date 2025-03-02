package github

import (
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackagePath(t *testing.T) {
	testCases := []struct {
		Reference string
		Expected  string
	}{
		{
			Reference: "ghcr.io/jmbannon/ytdl-sub",
			Expected:  "https://github.com/users/jmbannon/packages/container/package/ytdl-sub",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Reference, func(t *testing.T) {
			ref, err := oci.ParseReference(testCase.Reference)
			require.NoError(t, err)

			assert.Equal(t, testCase.Expected, PackageURL(ref))
		})
	}
}
