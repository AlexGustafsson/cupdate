package quay

type ScanStatus string

const (
	ScanStatusScanned     ScanStatus = "scanned"
	ScanStatusUnsupported            = "unsupported"
)

// Scan represents a scan result as defined by the Quay APIs.
type Scan struct {
	Status ScanStatus `json:"status"`
	Data   *struct {
		Layer struct {
			Name             string
			ParentName       string
			NamespaceName    string
			IndexedByVersion int
			Features         []Feature
		}
	} `json:"data"`
}

type Feature struct {
	Name            string
	VersionFormat   string
	NamespaceName   string
	AddedBy         string
	Version         string
	BaseScores      []float32
	CVEIds          []string
	Vulnerabilities []Vulnerability
}

type Severity string

const (
	SeverityCritical Severity = "Critical"
	SeverityHigh              = "High"
	SeverityMedium            = "Medium"
	SeverityLow               = "Low"
	SeverityUnknown           = "Unknown"
)

type Vulnerability struct {
	Severity      Severity
	NamespaceName string
	Link          string
	Description   string
	Name          string
	Metadata      VulnerabilityMetadata
}

type VulnerabilityMetadata struct {
	UpdatedBy     string
	RepoName      string
	RepoLink      string
	DistroName    string
	DistroVersion string
	NVD           map[string]any
}
