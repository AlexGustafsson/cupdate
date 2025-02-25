package oci

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/distribution/reference"
)

// Reference represents an OCI reference, i.e. an container image string.
type Reference struct {
	// Domain is the hostname of the registry.
	Domain string
	// Path is the namespace / project path of the reference.
	Path string

	// HasTag is true if the reference includes a tag.
	HasTag bool
	// Tag holds the tag specified in the reference.
	Tag string

	// HasDigest is true if the reference includes a digest.
	HasDigest bool
	// Digest holds the digest specified in the reference.
	Digest string
}

// Canonical converts the reference to its canonical form (i.e. with all fields
// explicitly set to their implicit default value).
// Panics if the refernece is invalid.
// A reference returned by [ParseReference] is always canonical.
func (r Reference) Canonical() Reference {
	// ParseReference always canonicalizes the reference, reuse it
	ref, err := ParseReference(r.String())
	if err != nil {
		panic(fmt.Sprintf("reference: canonical format misuse: %v", err))
	}

	return ref
}

// ParseReference parses a reference string.
// The returned reference is always canonical.
func ParseReference(v string) (Reference, error) {
	ref, err := reference.ParseNormalizedNamed(v)
	if err != nil {
		return Reference{}, err
	}

	hasTag := false
	tag := ""
	if tagged, ok := ref.(reference.Tagged); ok {
		hasTag = true
		tag = tagged.Tag()
	}

	hasDigest := false
	digest := ""
	if digested, ok := ref.(reference.Digested); ok {
		hasDigest = true
		digest = digested.Digest().String()
	}

	if !hasTag && !hasDigest {
		hasTag = true
		tag = "latest"
	}

	return Reference{
		Domain:    reference.Domain(ref),
		Path:      reference.Path(ref),
		HasTag:    hasTag,
		Tag:       tag,
		HasDigest: hasDigest,
		Digest:    digest,
	}, nil
}

// Name returns the name of the reference in a way typically used by users.
func (r Reference) Name() string {
	var builder strings.Builder

	if r.Domain == "docker.io" {
		r.Path = strings.TrimPrefix(r.Path, "library/")
		builder.WriteString(r.Path)
	} else if r.Domain == "" {
		builder.WriteString(r.Path)
	} else {
		builder.WriteString(r.Domain)
		builder.WriteString("/")
		builder.WriteString(r.Path)
	}

	return builder.String()
}

// Version is the familiar version of the reference, such as its tag, digest or
// "latest", if no tag or digest is specified.
// Mostly useful for human-readable use cases. For use with APIs, see
// [Reference.Reference].
func (r Reference) Version() string {
	if r.HasTag {
		return r.Tag
	} else if r.HasDigest {
		return r.Digest
	} else {
		return "latest"
	}
}

// Reference returns the reference as used by OCI distribution APIs.
func (r Reference) Reference() string {
	if r.HasDigest {
		return r.Digest
	} else if r.HasTag {
		return r.Tag
	} else {
		return "latest"
	}
}

// String returns the most compact way of describing the reference.
func (r Reference) String() string {
	var builder strings.Builder

	builder.WriteString(r.Name())

	if r.HasTag && (r.Tag != "latest" || r.HasDigest) {
		builder.WriteString(":")
		builder.WriteString(r.Tag)
	}

	if r.HasDigest {
		builder.WriteString("@")
		builder.WriteString(r.Digest)
	}

	return builder.String()
}

// MarshalJSON implements json.Marshaler.
func (r Reference) MarshalJSON() ([]byte, error) {
	v := r.String()
	return json.Marshal(v)
}

// MarshalJSON implements json.Unmarshaler.
func (r *Reference) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	ref, err := ParseReference(v)
	if err != nil {
		return err
	}

	*r = ref
	return nil
}
