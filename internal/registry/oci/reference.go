package oci

import (
	"encoding/json"

	"github.com/distribution/reference"
)

type Reference struct {
	Domain string
	Path   string

	HasTag bool
	Tag    string

	HasDigest bool
	Digest    string

	raw reference.Named
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

		raw: ref,
	}, nil
}

func (r Reference) Reference() reference.Named {
	return r.raw
}

func (r Reference) Name() string {
	return reference.FamiliarName(r.raw)
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
	return reference.FamiliarString(r.raw)
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
