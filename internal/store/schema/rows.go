package schema

import "time"

type RawImageRow struct {
	Reference     string     `sql:"reference,primary"`
	Tags          TagsBlob   `sql:"tags"`
	Graph         GraphBlob  `sql:"graph"`
	LastProcessed *time.Time `sql:"lastProcessed"`
}

type LinksRow struct {
	Reference string    `sql:"reference,primary"`
	Links     LinksBlob `sql:"links"`
}

type ReleaseNotesRow struct {
	Reference string    `sql:"reference,primary"`
	Title     string    `sql:"title"`
	HTML      string    `sql:"html"`
	Markdown  string    `sql:"markdown"`
	Released  time.Time `sql:"released"`
}

type DescriptionRow struct {
	Reference string `sql:"reference,primary"`
	HTML      string `sql:"html"`
	Markdown  string `sql:"markdown"`
}

type GraphRow struct {
	Reference string    `sql:"reference,primary"`
	Graph     GraphBlob `sql:"graph"`
}

type WorkflowRunRow struct {
	Reference string          `sql:"reference,primary"`
	Started   time.Time       `sql:"started"`
	Result    string          `sql:"result"`
	Blob      WorkflowRunBlob `sql:"blob"`
}

type ChangeRow struct {
	Reference string    `sql:"reference,primary"`
	Time      time.Time `sql:"time"`
	Type      string    `sql:"type"`

	ChangedBasic           bool `sql:"changedBasic"`
	ChangedLinks           bool `sql:"changedLinks"`
	ChangedReleaseNotes    bool `sql:"changedReleaseNotes"`
	ChangedDescription     bool `sql:"changedDescription"`
	ChangedGraph           bool `sql:"changedGraph"`
	ChangedVulnerabilities bool `sql:"changedVulnerabilities"`
}

type ScorecardRow struct {
	Reference string        `sql:"reference,primary"`
	Score     float64       `sql:"score"`
	Scorecard ScorecardBlob `sql:"scorecard"`
}

type TagRow struct {
	Reference string `sql:"reference"`
	Tag       string `sql:"tag"`
}

type UpdateRow struct {
	NewReference        string          `sql:"newReference,primary"`
	NewAnnotations      AnnotationsBlob `sql:"newAnnotations"`
	OldReference        string          `sql:"oldReference"`
	OldAnnotations      AnnotationsBlob `sql:"oldAnnotations"`
	VersionDiffSortable int             `sql:"versionDiffSortable"`
	Identified          time.Time       `sql:"identified"`
	Released            *time.Time      `sql:"released"`
}
