// Package models holds all models defined in the Cupdate API specification.
package models

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/osv"
)

type ImagePage struct {
	Images     []Image            `json:"images"`
	Summary    ImagePageSummary   `json:"summary"`
	Pagination PaginationMetadata `json:"pagination"`
}

type ImagePageSummary struct {
	Images     int `json:"images"`
	Outdated   int `json:"outdated"`
	Vulnerable int `json:"vulnerable"`
	Processing int `json:"processing"`
	Failed     int `json:"failed"`
}

type PaginationMetadata struct {
	Total int `json:"total"`
	// Page index. Starts at 1.
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Image struct {
	Reference           string      `json:"reference"`
	Created             *time.Time  `json:"created,omitempty"`
	LatestReference     string      `json:"latestReference,omitempty"`
	LatestCreated       *time.Time  `json:"latestCreated,omitempty"`
	VersionDiffSortable uint64      `json:"-"`
	Description         string      `json:"description,omitempty"`
	Tags                []string    `json:"tags"`
	Links               []ImageLink `json:"links"`
	Vulnerabilities     int         `json:"vulnerabilities"`
	LastModified        time.Time   `json:"lastModified"`
	Image               string      `json:"image,omitempty"`
}

type RawImage struct {
	Reference     string    `json:"reference"`
	Tags          []string  `json:"tags"`
	Graph         Graph     `json:"graph"`
	LastProcessed time.Time `json:"lastProcessed,omitempty"`
}

type ImageDescription struct {
	HTML     string `json:"html,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

type ImageReleaseNotes struct {
	Title    string    `json:"title"`
	HTML     string    `json:"html,omitempty"`
	Markdown string    `json:"markdown,omitempty"`
	Released time.Time `json:"released,omitempty"`
}

type ImageLink struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ImageVulnerability = osv.Vulnerability

type ImageScorecardRisk string

const (
	ImageScorecardRiskCritical ImageScorecardRisk = "critical"
	ImageScorecardRiskHigh     ImageScorecardRisk = "high"
	ImageScorecardRiskMedium   ImageScorecardRisk = "medium"
	ImageScorecardRiskLow      ImageScorecardRisk = "low"
)

type ImageScorecard struct {
	ReportURL  string             `json:"reportUrl"`
	Score      float32            `json:"score"`
	Risk       ImageScorecardRisk `json:"risk"`
	GenerateAt time.Time          `json:"generatedAt"`
}

type ImageProvenance struct {
	BuildInfo []ProvenanceBuildInfo `json:"buildInfo"`
}

type ProvenanceBuildInfo struct {
	ImageDigest         string            `json:"imageDigest"`
	Architecture        string            `json:"architecture,omitempty"`
	ArchitectureVariant string            `json:"architectureVariant,omitempty"`
	OperatingSystem     string            `json:"operatingSystem,omitempty"`
	Source              string            `json:"source,omitempty"`
	SourceRevision      string            `json:"sourceRevision,omitempty"`
	BuildStartedOn      time.Time         `json:"buildStartedOn,omitempty"`
	BuildFinishedOn     time.Time         `json:"buildFinishedOn,omitempty"`
	Dockerfile          string            `json:"dockerfile,omitempty"`
	BuildArguments      map[string]string `json:"buildArguments,omitempty"`
}

type ImageSBOM struct {
	SBOM []SBOM `json:"sbom"`
}

type SBOM struct {
	ImageDigest         string `json:"imageDigest"`
	Type                string `json:"type"`
	SBOM                string `json:"sbom"`
	Architecture        string `json:"architecture,omitempty"`
	ArchitectureVariant string `json:"architectureVariant,omitempty"`
	OperatingSystem     string `json:"operatingSystem,omitempty"`
}

type WorkflowRunResult string

const (
	WorkflowRunResultSucceeded WorkflowRunResult = "succeeded"
	WorkflowRunResultFailed    WorkflowRunResult = "failed"
)

type WorkflowRun struct {
	TraceID         string            `json:"traceId"`
	Started         time.Time         `json:"started"`
	DurationSeconds float64           `json:"duration"`
	Result          WorkflowRunResult `json:"result"`
	Jobs            []JobRun          `json:"jobs"`
}

type JobRunResult string

const (
	JobRunResultSucceeded JobRunResult = "succeeded"
	JobRunResultSkipped   JobRunResult = "skipped"
	JobRunResultFailed    JobRunResult = "failed"
)

type JobRun struct {
	Result          JobRunResult `json:"result"`
	Steps           []StepRun    `json:"steps"`
	DependsOn       []string     `json:"dependsOn"`
	JobID           string       `json:"jobId,omitempty"`
	JobName         string       `json:"jobName,omitempty"`
	Started         time.Time    `json:"started,omitempty"`
	DurationSeconds float64      `json:"duration,omitempty"`
}

type StepRunResult string

const (
	StepRunResultSucceeded StepRunResult = "succeeded"
	StepRunResultSkipped   StepRunResult = "skipped"
	StepRunResultFailed    StepRunResult = "failed"
)

type StepRun struct {
	Result          StepRunResult `json:"result"`
	StepName        string        `json:"stepName,omitempty"`
	Started         time.Time     `json:"started,omitempty"`
	DurationSeconds float64       `json:"duration,omitempty"`
	Error           string        `json:"error,omitempty"`
}

type Graph struct {
	Edges map[string]map[string]bool `json:"edges"`
	Nodes map[string]GraphNode       `json:"nodes"`
}

type GraphNode struct {
	Domain         string            `json:"domain"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels,omitempty"`
	InternalLabels map[string]string `json:"internalLabels,omitempty"`
}

type ImageEvent struct {
	Reference string    `json:"reference"`
	Type      EventType `json:"type"`
}

type EventType string

const (
	EventTypeImageUpdated             EventType = "imageUpdated"
	EventTypeImageProcessed           EventType = "imageProcessed"
	EventTypeImageNewVersionAvailable EventType = "imageNewVersionAvailable"
)
