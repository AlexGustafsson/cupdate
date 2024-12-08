package models

import "time"

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
}

type PaginationMetadata struct {
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

type Image struct {
	Reference       string               `json:"reference"`
	LatestReference string               `json:"latestReference,omitempty"`
	Description     string               `json:"description,omitempty"`
	Tags            []string             `json:"tags"`
	Links           []ImageLink          `json:"links"`
	Vulnerabilities []ImageVulnerability `json:"vulnerabilities"`
	LastModified    time.Time            `json:"lastModified"`
	Image           string               `json:"image,omitempty"`
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

type ImageVulnerability struct {
	ID          int    `json:"id"`
	Severity    string `json:"severity"`
	Authority   string `json:"authority"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

type Graph struct {
	Edges map[string]map[string]bool `json:"edges"`
	Nodes map[string]GraphNode       `json:"nodes"`
}

type GraphNode struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
	Name   string `json:"name"`
}
