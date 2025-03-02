package quay

type VulnerabilitySeverity string

const (
	VulnerabilitySeverityCritical    VulnerabilitySeverity = "critical"
	VulnerabilitySeverityHigh        VulnerabilitySeverity = "high"
	VulnerabilitySeverityMedium      VulnerabilitySeverity = "medium"
	VulnerabilitySeverityLow         VulnerabilitySeverity = "low"
	VulnerabilitySeverityUnspecified VulnerabilitySeverity = "unspecified"
)

type Vulnerability struct {
	Name        string
	Description string
	Links       []string
	// FeatureName is the name of the layer "feature", such as "json-c",
	// "libgomp".
	FeatureName    string
	FeatureVersion string
	// Layer is the digest of the layer that containers the vulnerability.
	Layer    string
	Severity VulnerabilitySeverity
}
