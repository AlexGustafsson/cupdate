package vulndb

import "time"

type Severity string

const (
	SeverityCritical    Severity = "critical"
	SeverityHigh        Severity = "high"
	SeverityMedium      Severity = "medium"
	SeverityLow         Severity = "low"
	SeverityUnspecified Severity = "unspecified"
)

type GitHubAdvisory struct {
	ID                string
	Repository        string
	Published         time.Time
	Severity          Severity
	IntroducedVersion string
	FixedVersion      string
}
