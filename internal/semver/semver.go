// Package semver contains methods of working with semantic versions.
package semver

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// NOTE: Parts of this code is translated from Renovate, which is under an
// AGPL-3.0 license.
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.ts
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/generic.ts#L18
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.spec.ts

var versionPattern = regexp.MustCompile(`^(?<version>\d+(?:\.\d+)*)(?<prerelease>\w*)$`)
var commitHashPattern = regexp.MustCompile(`^[a-f0-9]{7,40}$`)
var numericPattern = regexp.MustCompile(`^[0-9]+$`)

var (
	ErrUnsupportedVersionFormat = errors.New("unsupported version format")
)

type Version struct {
	Release    []int
	Suffix     string
	Prerelease string

	raw string
}

// IsStable returns whether or not the version is a "stable" version without a
// pre-release.
func (v *Version) IsStable() bool {
	return v.Prerelease == ""
}

// IsCompatible returns true if v can be compared to other.
func (v *Version) IsCompatible(other *Version) bool {
	return v.Suffix == other.Suffix && len(v.Release) == len(other.Release)
}

// Diff returns the type of bump that differs v to other.
// Returns an empty string in cases where a bump could not be found.
func (v *Version) Diff(other *Version) string {
	length := max(len(v.Release), len(other.Release))
	for i := 0; i < length; i++ {
		if other.Release[i] > v.Release[i] {
			switch i {
			case 0:
				return "major"
			case 1:
				return "minor"
			default:
				return "patch"
			}
		}
	}

	if other.Prerelease != v.Prerelease {
		return "patch"
	}

	return ""
}

// Compare compares two versions.
// Returns a negative number when v < other, a positive number when v > other
// and zero when v == other.
func (v *Version) Compare(other *Version) int {
	length := max(len(v.Release), len(other.Release))
	for i := 0; i < length; i++ {
		// Shorter is bigger 2.1 > 2.1.1
		if i >= len(v.Release) {
			return 1
		}
		if i >= len(other.Release) {
			return -1
		}

		part1 := v.Release[i]
		part2 := other.Release[i]
		if part1 != part2 {
			return part1 - part2
		}
	}

	if v.Prerelease != other.Prerelease {
		// Unstable is lower
		if v.Prerelease == "" && other.Prerelease != "" {
			return 1
		}
		if v.Prerelease != "" && other.Prerelease == "" {
			return -1
		}

		// Alphabetic order
		if v.Prerelease != "" && other.Prerelease != "" {
			return strings.Compare(v.Prerelease, other.Prerelease)
		}
	}

	// Equals
	return strings.Compare(other.Suffix, v.Suffix)
}

// Equals returns whether or not the two versions are equal.
func (v *Version) Equals(other *Version) bool {
	return v.Compare(other) == 0
}

// ParseVersion parses a [Version].
// Returns ErrUnsupportedVersionFormat if the version is not parsable.
func ParseVersion(version string) (*Version, error) {
	if version == "" {
		return nil, ErrUnsupportedVersionFormat
	}

	if commitHashPattern.MatchString(version) && !numericPattern.MatchString(version) {
		return nil, ErrUnsupportedVersionFormat
	}

	versionPieces := strings.Split(strings.TrimPrefix(version, "v"), "-")
	prefix := versionPieces[0]
	suffixPieces := versionPieces[1:]
	matchGroups := versionPattern.FindStringSubmatch(prefix)
	if matchGroups == nil {
		return nil, ErrUnsupportedVersionFormat
	}

	ver := matchGroups[1]
	prerelease := matchGroups[2]
	release := make([]int, 0)
	for _, x := range strings.Split(ver, ".") {
		n, err := strconv.ParseInt(x, 10, 32)
		if err != nil {
			return nil, err
		}
		release = append(release, int(n))
	}
	suffix := strings.Join(suffixPieces, "-")

	return &Version{
		Release:    release,
		Suffix:     suffix,
		Prerelease: prerelease,

		raw: version,
	}, nil
}
