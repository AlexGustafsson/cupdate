package oci

import (
	"encoding/json"
	"strings"

	"github.com/distribution/reference"
)

type Reference struct {
	Domain string
	Path   string

	HasTag bool
	Tag    string

	HasDigest bool
	Digest    string
}

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

	return Reference{
		Domain:    reference.Domain(ref),
		Path:      reference.Path(ref),
		HasTag:    hasTag,
		Tag:       tag,
		HasDigest: hasDigest,
		Digest:    digest,
	}, nil
}

func (r Reference) Name() string {
	var builder strings.Builder

	if r.Domain == "docker.io" {
		r.Path = strings.TrimPrefix(r.Path, "library/")
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
func (r Reference) Version() string {
	if r.HasTag {
		return r.Tag
	} else if r.HasDigest {
		return r.Digest
	} else {
		return "latest"
	}
}

func (r Reference) String() string {
	var builder strings.Builder

	builder.WriteString(r.Name())

	if r.HasTag {
		if r.Tag != "latest" {
			builder.WriteString(":")
			builder.WriteString(r.Tag)
		}
	} else if r.HasDigest {
		builder.WriteString("@")
		builder.WriteString(r.Digest)
	}

	return builder.String()
}

func (r Reference) MarshalJSON() ([]byte, error) {
	v := r.String()
	return json.Marshal(v)
}

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
