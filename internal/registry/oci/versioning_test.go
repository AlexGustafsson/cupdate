package oci

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: Some test cases are from Renovate, which is under an
// AGPL-3.0 license.
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.ts
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/generic.ts#L18
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.spec.ts

func TestVersionIsCompatible(t *testing.T) {
	testCases := []struct {
		Version  string
		Range    string
		Expected bool
	}{
		{
			Version:  "3.7.0",
			Range:    "3.7.0",
			Expected: true,
		},
		{
			Version:  "3.7.0b1",
			Range:    "3.7.0",
			Expected: true,
		},
		{
			Version:  "3.7-alpine",
			Range:    "3.7.0",
			Expected: false,
		},
		{
			Version:  "3.8.0-alpine",
			Range:    "3.7.0",
			Expected: false,
		},
		{
			Version:  "3.8.0b1-alpine",
			Range:    "3.7.0",
			Expected: false,
		},
		{
			Version:  "3.8.2",
			Range:    "3.7.0",
			Expected: true,
		},
		{
			Version:  "3.7.0",
			Range:    "3.7.0-alpine",
			Expected: false,
		},
		{
			Version:  "3.7.0b1",
			Range:    "3.7.0-alpine",
			Expected: false,
		},
		{
			Version:  "3.7-alpine",
			Range:    "3.7.0-alpine",
			Expected: false,
		},
		{
			Version:  "3.8.0-alpine",
			Range:    "3.7.0-alpine",
			Expected: true,
		},
		{
			Version:  "3.8.0b1-alpine",
			Range:    "3.7.0-alpine",
			Expected: true,
		},
		{
			Version:  "3.8.2",
			Range:    "3.7.0-alpine",
			Expected: false,
		},
		{
			Version:  "0.7.2",
			Range:    "0.8.5rc51",
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s, %s: %v", testCase.Version, testCase.Range, testCase.Expected), func(t *testing.T) {
			v, err := ParseVersion(testCase.Version)
			require.NoError(t, err)

			r, err := ParseVersion(testCase.Range)
			require.NoError(t, err)

			assert.Equal(t, testCase.Expected, v.IsCompatible(r))
		})
	}
}
func TestVersionIsStable(t *testing.T) {
	testCases := []struct {
		Version  string
		Expected bool
	}{
		{
			Version:  "3.7.0",
			Expected: true,
		},
		{
			Version:  "3.7.0b1",
			Expected: false,
		},
		{
			Version:  "3.7-alpine",
			Expected: true,
		},
		{
			Version:  "3.8.0-alpine",
			Expected: true,
		},
		{
			Version:  "3.8.0b1-alpine",
			Expected: false,
		},
		{
			Version:  "3.8.2",
			Expected: true,
		},
		{
			Version:  "0.8.5rc51",
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s: %v", testCase.Version, testCase.Expected), func(t *testing.T) {
			v, err := ParseVersion(testCase.Version)
			require.NoError(t, err)

			assert.Equal(t, testCase.Expected, v.IsStable())
		})
	}
}

func TestCompareVersionsUnstable(t *testing.T) {
	versions := []string{
		"3.7.0",
		"3.7-alpine",
		"3.7.0b1",
		"3.7.0b5",
		"3.8.3rc1",
		"3.8.0b1-alpine",
		"3.8.0-alpine",
		"3.8.2",
		"3.8.0",
		"3.8.1rc1",
	}

	expected := []string{
		"3.7.0b1",
		"3.7.0b5",
		"3.7.0",
		"3.7-alpine",
		"3.8.0b1-alpine",
		"3.8.0-alpine",
		"3.8.0",
		"3.8.1rc1",
		"3.8.2",
		"3.8.3rc1",
	}

	slices.SortFunc(versions, func(as, bs string) int {
		// Parse whilst sorting so that we get a nice diff based on the strings
		// above rather than parsed versions
		a, err := ParseVersion(as)
		require.NoError(t, err)

		b, err := ParseVersion(bs)
		require.NoError(t, err)

		return CompareVersions(a, b)
	})

	assert.Equal(t, versions, expected)
}

func TestCompareVersionsSemver(t *testing.T) {
	versions := []string{
		"1.1.1",
		"1.2.3",
		"2.0.1",
		"1.2.3",
		"1.2.3",
		"1.3.4",
		"1.2.3",
		"0.9.5",
	}

	expected := []string{
		"0.9.5",
		"1.1.1",
		"1.2.3",
		"1.2.3",
		"1.2.3",
		"1.2.3",
		"1.3.4",
		"2.0.1",
	}

	slices.SortFunc(versions, func(as, bs string) int {
		// Parse whilst sorting so that we get a nice diff based on the strings
		// above rather than parsed versions
		a, err := ParseVersion(as)
		require.NoError(t, err)

		b, err := ParseVersion(bs)
		require.NoError(t, err)

		return CompareVersions(a, b)
	})

	assert.Equal(t, versions, expected)
}
