package schema

import "time"

type RawImage struct {
	Reference     string         `sql:"reference,primary"`
	Tags          Blob[[]string] `sql:"tags"`
	Graph         Blob[Graph]    `sql:"graph"`
	LastProcessed *time.Time     `sql:"lastProcessed"`
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

type Links struct {
	Reference string         `sql:"reference,primary"`
	Links     Blob[[]string] `sql:"links"`
}

type ReleaseNote struct {
	Reference string    `sql:"reference,primary"`
	Title     string    `sql:"title"`
	HTML      string    `sql:"html"`
	Markdown  string    `sql:"markdown"`
	Released  time.Time `sql:"released"`
}

type Description struct {
	Reference string `sql:"reference,primary"`
	HTML      string `sql:"html"`
	Markdown  string `sql:"markdown"`
}

type Tag struct {
	Reference string `sql:"reference"`
	Tag       string `sql:"tag"`
}

type Update struct {
	NewReference        string                  `sql:"newReference,primary"`
	NewAnnotations      Blob[map[string]string] `sql:"newAnnotations"`
	OldReference        string                  `sql:"oldReference"`
	OldAnnotations      Blob[map[string]string] `sql:"oldAnnotations"`
	VersionDiffSortable int                     `sql:"versionDiffSortable"`
	Identified          time.Time               `sql:"identified"`
	Released            *time.Time              `sql:"released"`
}
