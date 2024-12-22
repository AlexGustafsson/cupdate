package semver

// CompareVersions compares two compatible [Version]s.
// Returns a negative number when a < b, a positive number when a > b and zero
// when a == b.
// Useful to use as a sort func.
func CompareVersions(a *Version, b *Version) int {
	return a.Compare(b)
}

// LatestOpinionatedVersionString returns the latest compatible version and
// whether or not a compatible version was found.
// If there is a new version on the same major as the current, it is preferred.
// If the current version is the newest version of the major, the latest version
// is returned. This allows more graceful handling of version tracks for things
// like databases where multiple majors are supported concurrently.
func LatestOpinionatedVersionString(current string, versions []string) (string, bool) {
	if current == "latest" {
		return current, true
	}

	currentVersion, err := ParseVersion(current)
	if err != nil {
		// Despite this being a parse error, handle it as a compatible version in
		// context of the current version
		return current, true
	}

	compatibleVersions := make([]*Version, 0)
	for _, version := range versions {
		v, err := ParseVersion(version)
		if err == nil && v.IsCompatible(currentVersion) {
			// TODO: Decide on how to handle semantic versions - should we be strict?
			// None seem to be using correct semver release formats, so I guess we
			// need to parse losely. Right now the correct semver release candidate
			// "-rc" is parsed as a suffix. The invalid release candidate "8.0.3rc1"
			// is parsed as  a release candidate...
			// Perhaps conduct some data-driven experiment and check the state of the
			// ecosystem?
			if v.Prerelease == "" || v.Prerelease != "" && currentVersion.Prerelease != "" {
				compatibleVersions = append(compatibleVersions, v)
			}
		}
	}

	if len(compatibleVersions) == 0 {
		return current, false
	}

	latestVersionOfSameMajor := currentVersion
	latestVersion := currentVersion
	compatibleFound := false
	for _, version := range compatibleVersions {
		if len(version.Release) > 0 && len(currentVersion.Release) > 0 &&
			version.Release[0] == currentVersion.Release[0] &&
			version.Compare(latestVersionOfSameMajor) > 0 {
			latestVersionOfSameMajor = version
		}

		ret := version.Compare(latestVersion)
		if ret > 0 {
			latestVersion = version
		}
		if ret >= 0 {
			compatibleFound = true
		}
	}

	// End of major track, return the latest version
	if latestVersionOfSameMajor.Equals(currentVersion) {
		return latestVersion.raw, compatibleFound
	}

	// The current version is newer than any version in the versions list
	if latestVersionOfSameMajor.Compare(currentVersion) < 0 {
		return currentVersion.raw, true
	}

	return latestVersionOfSameMajor.raw, true
}

// PackInt64 packs a [Version] into a lossy format which fits into a 64-bit int.
// The resulting int is sortable when compared to compatible versions. That is,
// a version of a higher major, higher minor, higher patch and so on will be a
// higher value than a lower version. Useful for creating a value which is then
// used to diff two versions of the same image. The resulting diff is sortable
// among with calculated values. Might fall apart if there are many parts, but
// most of the time, there are only three or four version parts.
func PackInt64(version *Version) uint64 {
	var packed uint64

	bitsPerPart := 64 / len(version.Release)
	for i, part := range version.Release {
		packed |= uint64(part) << uint64((len(version.Release)-i-1)*bitsPerPart)
	}

	return packed
}
