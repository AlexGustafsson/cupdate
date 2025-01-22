package oci

import (
	"strings"
	"time"
)

type Annotations map[string]string

func (a Annotations) Source() string {
	if a == nil {
		return ""
	}

	s := a["org.opencontainers.image.source"]
	if s == "" {
		s = a["org.label-schema.vcs-url"]
	}

	return s
}

func (a Annotations) CreatedTime() time.Time {
	s := a["org.opencontainers.image.created"]
	if s == "" {
		s = a["org.label-schema.build-date"]
	}

	// Golang's RFC3339 format requires a "T", but the spec allows a space
	s = strings.Replace(s, " ", "T", 1)
	time, _ := time.Parse(time.RFC3339, s)
	return time
}

func (a Annotations) URL() string {
	return a["org.opencontainers.image.url"]
}

func (a Annotations) DocumentationURL() string {
	return a["org.opencontainers.image.documentation"]
}
