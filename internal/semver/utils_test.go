package semver

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareVersions(t *testing.T) {
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

func TestLatestVersionOnSameTrack(t *testing.T) {
	versions := []string{
		"5",
		"5-focal",
		"5-nanoserver",
		"5-nanoserver-1809",
		"5-nanoserver-ltsc2022",
		"5-windowsservercore-1809",
		"5-windowsservercore-ltsc2022",
		"5.0.30",
		"5.0.30-focal",
		"5.0.30-nanoserver",
		"5.0.30-nanoserver-1809",
		"5.0.30-nanoserver-ltsc2022",
		"5.0.30-windowsservercore",
		"5.0.30-windowsservercore-1809",
		"5.0.30-windowsservercore-ltsc2022",
		"6",
		"6-jammy",
		"6-nanoserver",
		"6-nanoserver-1809",
		"6-nanoserver-ltsc2022",
		"6-windowsservercore-1809",
		"6-windowsservercore-ltsc2022",
		"6.0-rc",
		"6.0-rc-jammy",
		"6.0-rc-nanoserver",
		"6.0-rc-nanoserver-1809",
		"6.0-rc-nanoserver-ltsc2022",
		"6.0-rc-windowsservercore-1809",
		"6.0-rc-windowsservercore-ltsc2022",
		"6.0.19",
		"6.0.19-jammy",
		"6.0.19-nanoserver",
		"6.0.19-nanoserver-1809",
		"6.0.19-nanoserver-ltsc2022",
		"6.0.19-windowsservercore",
		"6.0.19-windowsservercore-1809",
		"6.0.19-windowsservercore-ltsc2022",
		"6.0.20-rc3",
		"6.0.20-rc3-jammy",
		"6.0.20-rc3-nanoserver",
		"6.0.20-rc3-nanoserver-1809",
		"6.0.20-rc3-nanoserver-ltsc2022",
		"6.0.20-rc3-windowsservercore",
		"6.0.20-rc3-windowsservercore-1809",
		"6.0.20-rc3-windowsservercore-ltsc2022",
		"7",
		"7-jammy",
		"7-nanoserver",
		"7-nanoserver-1809",
		"7-nanoserver-ltsc2022",
		"7-windowsservercore-1809",
		"7-windowsservercore-ltsc2022",
		"7.0-rc",
		"7.0-rc-jammy",
		"7.0-rc-nanoserver",
		"7.0-rc-nanoserver-1809",
		"7.0-rc-nanoserver-ltsc2022",
		"7.0-rc-windowsservercore-1809",
		"7.0-rc-windowsservercore-ltsc2022",
		"7.0.15",
		"7.0.15-jammy",
		"7.0.15-nanoserver",
		"7.0.15-nanoserver-1809",
		"7.0.15-nanoserver-ltsc2022",
		"7.0.15-windowsservercore",
		"7.0.15-windowsservercore-1809",
		"7.0.15-windowsservercore-ltsc2022",
		"7.0.16-rc1",
		"7.0.16-rc1-jammy",
		"7.0.16-rc1-nanoserver",
		"7.0.16-rc1-nanoserver-1809",
		"7.0.16-rc1-nanoserver-ltsc2022",
		"7.0.16-rc1-windowsservercore",
		"7.0.16-rc1-windowsservercore-1809",
		"7.0.16-rc1-windowsservercore-ltsc2022",
		"8.0.4",
		"8.0.4-nanoserver",
		"8.0.4-nanoserver-1809",
		"8.0.4-nanoserver-ltsc2022",
		"8.0.4-noble",
		"8.0.4-windowsservercore",
		"8.0.4-windowsservercore-1809",
		"8.0.4-windowsservercore-ltsc2022",
		"nanoserver",
		"nanoserver-1809",
		"nanoserver-ltsc2022",
		"noble",
		"windowsservercore-1809",
		"windowsservercore-ltsc2022",
	}

	testCases := []struct {
		Version  string
		Expected string
		OK       bool
	}{
		{
			// End of major, expected latest major
			Version:  "5",
			Expected: "7",
			OK:       true,
		},
		{
			// Patch on same major track
			Version:  "5.0.1",
			Expected: "5.0.30",
			OK:       true,
		},
		{
			// Does not exist, assume newer
			Version:  "8.1.0",
			Expected: "8.1.0",
			OK:       false,
		},
		{
			// Latest convention
			Version:  "latest",
			Expected: "latest",
			OK:       true,
		},
		{
			// Existing tag with no apparent semver
			Version:  "noble",
			Expected: "noble",
			OK:       true,
		},
		{
			// Non-existing tag with no apparent semver
			Version:  "royal",
			Expected: "royal",
			OK:       true,
		},
		{
			// Version with specified distro
			Version:  "6.0.18-jammy",
			Expected: "6.0.19-jammy",
			OK:       true,
		},
		{
			// Release candidate bumps are not supported as they are not well defined,
			// handle like a patch with a specified distro (in this case, end of a
			// major)
			Version:  "6.0.20-rc2",
			Expected: "6.0.20-rc2",
			OK:       false,
		},
		{
			// Apparent end of same major track - recommend newest
			Version:  "6.0.19",
			Expected: "8.0.4",
			OK:       true,
		},
		{
			// Latest available
			Version:  "8.0.4",
			Expected: "8.0.4",
			OK:       true,
		},
		// TODO: Decide on how to handle semantic versions - should we be strict?
		// None seem to be using correct semver release formats, so I guess we need
		// to parse losely. Right now the correct semver release candidate "-rc" is
		// parsed as a suffix. The invalid release candidate "8.0.3rc1" is parsed as
		// a release candidate...
		// Perhaps conduct some data-driven experiment and check the state of the
		// ecosystem?
		// {
		// 	// From pre-release to non-pre-release is fine
		// 	Version:  "8.0.3-rc1",
		// 	Expected: "8.0.4",
		// 	OK:       true,
		// },
		// {
		// 	// From non-pre-release to pre-release is not fine
		// 	Version:  "8.0.3",
		// 	Expected: "8.0.4-rc1",
		// 	OK:       true,
		// },
	}

	for _, testCase := range testCases {
		t.Run(testCase.Version, func(t *testing.T) {
			actual, ok := LatestOpinionatedVersionString(testCase.Version, versions)
			assert.Equal(t, testCase.Expected, actual)
			assert.Equal(t, testCase.OK, ok)
		})
	}
}

func TestPackInt64(t *testing.T) {
	// Assert that the zero value is a non-zero value (to differ from an unset
	// value)
	assert.Equal(t, uint64(1), PackInt64(nil))

	// Assert that 1 patch diff is 1, skipping the zero value
	assert.Equal(t, uint64(2), PackedSingleDigitPatchDiff)

	expected := []string{
		"0",
		"0.0",
		"0.0.0",
		"0.0.1",
		"0.0.2",
		"0.2.0",
		"1.0.0",
		"1.0.0-alpine",
		"1.2.3",
		"1.2.4",
	}

	actual := append([]string{}, expected...)

	rand.New(rand.NewSource(5325)).Shuffle(len(actual), func(i, j int) {
		actual[i], actual[j] = actual[j], actual[i]
	})

	slices.SortStableFunc(actual, func(a string, b string) int {
		av, err := ParseVersion(a)
		require.NoError(t, err)

		bv, err := ParseVersion(b)
		require.NoError(t, err)

		return int(PackInt64(av) - PackInt64(bv))
	})

	assert.Equal(t, expected, actual)
}
