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

// Merge returns the merge of a and b.
// If both are nil, nil is returned.
func (a Annotations) Merge(b Annotations) Annotations {
	if a == nil && b == nil {
		return nil
	} else if a == nil {
		return b
	} else if b == nil {
		return b
	} else {
		clone := make(Annotations)
		for k, v := range a {
			clone[k] = v
		}
		for k, v := range b {
			clone[k] = v
		}
		return clone
	}
}
