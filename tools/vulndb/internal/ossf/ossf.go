package ossf

import "time"

// A schema for describing a vulnerability in an open source package. See also
// https://ossf.github.io/osv-schema/.
type OpenSourceVulnerability struct {
	Affected         []Affected     `json:"affected,omitempty"`
	Aliases          []string       `json:"aliases,omitempty"`
	Credits          []Credit       `json:"credits,omitempty"`
	DatabaseSpecific map[string]any `json:"database_specific,omitempty"`
	Details          *string        `json:"details,omitempty"`
	ID               string         `json:"id"`
	Modified         time.Time      `json:"modified"`
	Published        time.Time      `json:"published,omitempty"`
	References       []Reference    `json:"references,omitempty"`
	Related          []string       `json:"related,omitempty"`
	SchemaVersion    *string        `json:"schema_version,omitempty"`
	Severity         Severity       `json:"severity,omitempty"`
	Summary          *string        `json:"summary,omitempty"`
	Withdrawn        time.Time      `json:"withdrawn,omitempty"`
}

type Affected struct {
	DatabaseSpecific  map[string]any   `json:"database_specific,omitempty"`
	EcosystemSpecific map[string]any   `json:"ecosystem_specific,omitempty"`
	Package           *AffectedPackage `json:"package,omitempty"`
	Ranges            []AffactedRange  `json:"ranges,omitempty"`
	Severity          Severity         `json:"severity,omitempty"`
	Versions          []string         `json:"versions,omitempty"`
}

type AffectedPackage struct {
	Ecosystem string  `json:"ecosystem"`
	Name      string  `json:"name"`
	Purl      *string `json:"purl,omitempty"`
}

type AffactedRange struct {
	DatabaseSpecific map[string]any   `json:"database_specific,omitempty"`
	Events           []map[string]any `json:"events"`
	Repo             *string          `json:"repo,omitempty"`
	Type             string           `json:"type"`
}

type Credit struct {
	Contact []string `json:"contact,omitempty"`
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
}

type Reference struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Severity []struct {
	Score string `json:"score"`
	Type  string `json:"type"`
}
