// Package semver contains methods of working with semantic versions.
//
// NOTE: Parts of this code is translated from Renovate, which is under an
// AGPL-3.0 license.
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.ts
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/generic.ts#L18
//
// SEE: https://github.com/renovatebot/renovate/blob/4a9b489b71f19443c352cd5ae045d93264204120/lib/modules/versioning/docker/index.spec.ts
package semver

import (
	"regexp"
	"strconv"
	"strings"
)

var versionPattern = regexp.MustCompile(`^(?<version>\d+(?:\.\d+)*)(?<prerelease>\w*)$`)
var commitHashPattern = regexp.MustCompile(`^[a-f0-9]{7,40}$`)
var numericPattern = regexp.MustCompile(`^[0-9]+$`)

type Version struct {
	Release    []int
	Suffix     string
	Prerelease string
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

// Compare returns 1 if other is newer than v, -1 if v is newer than other and
// 0 if the two versions are equal.
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

// ParseVersion parses a [Version].
func ParseVersion(version string) (*Version, error) {
	if commitHashPattern.MatchString(version) && !numericPattern.MatchString(version) {
		return nil, nil
	}

	versionPieces := strings.Split(strings.TrimPrefix(version, "v"), "-")
	prefix := versionPieces[0]
	suffixPieces := versionPieces[1:]
	matchGroups := versionPattern.FindStringSubmatch(prefix)
	if matchGroups == nil {
		return nil, nil
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
	}, nil
}
