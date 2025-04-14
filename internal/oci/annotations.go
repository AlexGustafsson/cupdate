package oci

import (
	"strings"
	"time"
)

// Annotations holds OCI annotations.
type Annotations map[string]string

// NOTE: Annotations below prioritize standardize annotations and are listed in
// the order of the defining document.
// SEE: https://github.com/opencontainers/image-spec/blob/c05acf7eb327dae4704a4efe01253a0e60af6b34/annotations.md#annotations
// SEE: https://docs.redhat.com/en/documentation/red_hat_openstack_platform/13/html/partner_integration/building-certified-container-images#dockerfile-requirements

// Created is the date and time on which the image was built.
func (a Annotations) Created() time.Time {
	if a == nil {
		return time.Time{}
	}

	s := a.oneOf(
		"org.opencontainers.image.created",
		"org.label-schema.build-date",
	)

	// Golang's RFC3339 format requires a "T", but the spec allows a space.
	s = strings.Replace(s, " ", "T", 1)
	time, _ := time.Parse(time.RFC3339, s)
	return time
}

// Authors contains contact details of the people or organization responsible
// for the image.
func (a Annotations) Authors() string {
	return a.oneOf(
		"org.opencontainers.image.authors",
		"maintainer",
	)
}

// URL is a URL to find more information on the image.
func (a Annotations) URL() string {
	return a.oneOf(
		"org.opencontainers.image.url",
		"org.label-schema.url",
	)
}

// Documentation is a URL to get documentation on the image.
func (a Annotations) Documentation() string {
	return a.oneOf(
		"org.opencontainers.image.documentation",
		"org.label-schema.usage",
	)
}

// Source is a URL to get source code for building the image.
func (a Annotations) Source() string {
	return a.oneOf(
		"org.opencontainers.image.source",
		"org.label-schema.vcs-url",
	)
}

// Version is the version of the packaged software.
func (a Annotations) Version() string {
	return a.oneOf(
		"org.opencontainers.image.version",
		"org.label-schema.version",
		"version",
	)
}

// Revision is the source control revision identifier for the packaged software.
func (a Annotations) Revision() string {
	return a.oneOf(
		"org.opencontainers.image.revision",
		"org.label-schema.vcs-ref",
	)
}

// Vendor is the name of the distributing entity, organization or individual.
func (a Annotations) Vendor() string {
	return a.oneOf(
		"org.opencontainers.image.vendor",
		"org.label-schema.vendor",
		"vendor",
	)
}

// Licenses is the License(s) under which contained software is distributed, as
// an SPDX License Expression.
func (a Annotations) Licenses() string {
	return a.oneOf(
		"org.opencontainers.image.licenses",
	)
}

// RefName is the name of the reference for a target (string).
func (a Annotations) RefName() string {
	return a.oneOf(
		"org.opencontainers.image.ref.name",
	)
}

// Title is a human-readable title of the image.
func (a Annotations) Title() string {
	return a.oneOf(
		"org.opencontainers.image.title",
		"org.label-schema.name",
		"summary",
	)
}

// Description is a human-readable description of the software packaged in the
// image.
func (a Annotations) Description() string {
	return a.oneOf(
		"org.opencontainers.image.description",
		"org.label-schema.description",
		"description",
	)
}

// BaseDigest is the digest of the image this image is based on.
func (a Annotations) BaseDigest() string {
	return a.oneOf(
		"org.opencontainers.image.base.digest",
	)
}

// BaseName is the digest of the image this image is based on.
func (a Annotations) BaseName() string {
	return a.oneOf(
		"org.opencontainers.image.base.name",
	)
}

// DockerReferenceType describes the type of artifact.
// Used for attestation manifests.
func (a Annotations) DockerReferenceType() string {
	return a.oneOf("vnd.docker.reference.type")
}

// DockerReferenceDigest is the digest of the image for which the reference type
// annotation is valid.
// Used for attestation manifests.
func (a Annotations) DockerReferenceDigest() string {
	return a.oneOf("vnd.docker.reference.digest")
}

// InTotoPredicateType returns the predicate type of a layer.
// SEE: https://in-toto.io.
// SEE: https://github.com/in-toto/attestation/tree/v1.0/spec/predicates.
func (a Annotations) InTotoPredicateType() string {
	return a.oneOf("in-toto.io/predicate-type")
}

func (a Annotations) oneOf(keys ...string) string {
	for _, key := range keys {
		if value := a[key]; value != "" {
			return value
		}
	}

	return ""
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
