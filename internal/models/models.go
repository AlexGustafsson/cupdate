package models

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type ImagePage struct {
	Images     []*Image            `json:"images"`
	Summary    *ImagePageSummary   `json:"summary"`
	Pagination *PaginationMetadata `json:"pagination"`
}

type ImagePageSummary struct {
	Images   int `json:"images,omitempty"`
	Outdated int `json:"outdated,omitempty"`
	Pods     int `json:"pods,omitempty"`
}

type PaginationMetadata struct {
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Image struct {
	Name           string       `json:"name"`
	Description    string       `json:"description,omitempty"`
	CurrentVersion string       `json:"currentVersion"`
	LatestVersion  string       `json:"latestVersion"`
	Tags           []string     `json:"tags"`
	Links          []*ImageLink `json:"links"`
	Image          string       `json:"image,omitempty"`
}

type ImageDescription struct {
	HTML     string `json:"html,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

type ImageReleaseNotes struct {
	Title string `json:"title"`
	HTML  string `json:"html,omitempty"`
}

type ImageLink struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Graph struct {
	Root *GraphNode `json:"root"`
}

type GraphNode struct {
	Domain  string       `json:"domain"`
	Type    string       `json:"type"`
	Name    string       `json:"name"`
	Parents []*GraphNode `json:"parents"`
}
