package semver

// CompareVersions compares [Version]s.
// Useful to use as a sort func.
func CompareVersions(a *Version, b *Version) int {
	return a.Compare(b)
}
