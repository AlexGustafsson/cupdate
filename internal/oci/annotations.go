package oci

import (
	"strings"
	"time"
)

// Annotations holds OCI annotations.
type Annotations map[string]string

// Source returns the value of the standard OCI source annotations.
// The annotation should point to a source code repository.
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

// CreatedTime returns value of the standard OCI created time annotations.
// The annotation should hold the time at which the image or artifact was
// created.
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

// URL returns value of the standard OCI URL annotation.
// The annotation should hold a URL to the project's website.
func (a Annotations) URL() string {
	return a["org.opencontainers.image.url"]
}

// DocumentationURL returns value of the standard OCI documentation URL
// annotation.
// The annotation should hold a URL to the project's documentation.
func (a Annotations) DocumentationURL() string {
	return a["org.opencontainers.image.documentation"]
}

// Merge returns the merge of a and b.
// If both are nil, nil is returned.
// If the values exist in both sets, values in b takes precedence.
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
