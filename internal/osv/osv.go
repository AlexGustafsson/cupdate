package osv

import (
	"slices"
	"time"
)

// A schema for describing a vulnerability in an open source package. See also
// https://ossf.github.io/osv-schema/.
type Vulnerability struct {
	ID               string         `json:"id"`
	Modified         time.Time      `json:"modified"`
	Affected         []Affected     `json:"affected,omitempty"`
	Aliases          []string       `json:"aliases,omitempty"`
	Credits          []Credit       `json:"credits,omitempty"`
	DatabaseSpecific map[string]any `json:"database_specific,omitempty"`
	Details          string         `json:"details,omitempty"`
	Published        *time.Time     `json:"published,omitempty"`
	References       []Reference    `json:"references,omitempty"`
	Related          []string       `json:"related,omitempty"`
	SchemaVersion    string         `json:"schema_version,omitempty"`
	Severities       []Severity     `json:"severity,omitempty"`
	Summary          string         `json:"summary,omitempty"`
	Withdrawn        *time.Time     `json:"withdrawn,omitempty"`
}

type Affected struct {
	DatabaseSpecific  map[string]any   `json:"database_specific,omitempty"`
	EcosystemSpecific map[string]any   `json:"ecosystem_specific,omitempty"`
	Package           *AffectedPackage `json:"package,omitempty"`
	Ranges            []AffectedRange  `json:"ranges,omitempty"`
	Severities        []Severity       `json:"severity,omitempty"`
	Versions          []string         `json:"versions,omitempty"`
}

type AffectedPackage struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Purl      string `json:"purl,omitempty"`
}

type AffectedRange struct {
	Type             string         `json:"type"`
	DatabaseSpecific map[string]any `json:"database_specific,omitempty"`
	Events           []Event        `json:"events,omitempty"`
	Repo             string         `json:"repo,omitempty"`
}

type Event struct {
	Introduced   string `json:"introduced,omitempty"`
	Fixed        string `json:"fixed,omitempty"`
	LastAffected string `json:"last_affected,omitempty"`
	Limit        string `json:"limit,omitempty"`
}

type Credit struct {
	Name    string   `json:"name"`
	Contact []string `json:"contact,omitempty"`
	Type    string   `json:"type,omitempty"`
}

type ReferenceType string

const (
	ReferenceTypeAdvisory   ReferenceType = "ADVISORY"
	ReferenceTypeArticle    ReferenceType = "ARTICLE"
	ReferenceTypeDetection  ReferenceType = "DETECTION"
	ReferenceTypeDiscussion ReferenceType = "DISCUSSION"
	ReferenceTypeReport     ReferenceType = "REPORT"
	ReferenceTypeFix        ReferenceType = "FIX"
	ReferenceTypeGit        ReferenceType = "GIT"
	ReferenceTypeIntroduced ReferenceType = "INTRODUCED"
	ReferenceTypePackage    ReferenceType = "PACKAGE"
	ReferenceTypeEvidence   ReferenceType = "EVIDENCE"
	ReferenceTypeWeb        ReferenceType = "WEB"
)

type Reference struct {
	Type ReferenceType `json:"type"`
	URL  string        `json:"url"`
}

type Severity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

type NormalizedSeverity string

const (
	NormalizedSeverityCritical    NormalizedSeverity = "critical"
	NormalizedSeverityHigh        NormalizedSeverity = "high"
	NormalizedSeverityMedium      NormalizedSeverity = "medium"
	NormalizedSeverityLow         NormalizedSeverity = "low"
	NormalizedSeverityUnspecified NormalizedSeverity = "unspecified"
)

func (o Vulnerability) NormalizedSeverity() NormalizedSeverity {
	if o.DatabaseSpecific != nil {
		switch o.DatabaseSpecific["severity"] {
		case "CRITICAL":
			return NormalizedSeverityCritical
		case "HIGH":
			return NormalizedSeverityHigh
		case "MODERATE", "MEDIUM":
			return NormalizedSeverityMedium
		case "LOW":
			return NormalizedSeverityLow
		}
	}

	// TODO: Get the severities from packages / parse CVSS score?
	return NormalizedSeverityUnspecified
}

// Compare returns negative if s is of higher severity than o. Positive if s is
// of lower severity than o. 0 if equivalent (or both severities are unknown).
//
// Useful to sort severities using functions like [slices.Sort].
func (s NormalizedSeverity) Compare(o NormalizedSeverity) int {
	if s == o {
		return 0
	}

	order := []NormalizedSeverity{
		NormalizedSeverityCritical,
		NormalizedSeverityHigh,
		NormalizedSeverityMedium,
		NormalizedSeverityLow,
		NormalizedSeverityUnspecified,
	}

	orderS := slices.Index(order, s)
	if orderS == -1 {
		orderS = len(order)
	}

	orderO := slices.Index(order, o)
	if orderO == -1 {
		orderO = len(order)
	}

	return orderS - orderO
}
