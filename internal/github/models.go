package github

import "time"

type Release struct {
	Owner       string
	Repository  string
	Tag         string
	Released    time.Time
	Title       string
	Description string
	URL         string
}

type Package struct {
	Owner      string
	Repository string
	Tags       []PackageTag
}

type PackageTag struct {
	Name   string
	Latest bool
}
