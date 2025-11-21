package platform

import (
	"fmt"
	"maps"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/semver"
)

// Labels holds labels / annotations found by platform implementations, which
// map to things like Docker labels or Kubernetes resource annotations.
type Labels map[string]string

// Ignore returns true if the Cupdate ignore label is set to true.
func (l Labels) Ignore() bool {
	v, _ := l.oneOf("config.cupdate/ignore", "cupdate.config.ignore")
	return v == "true"
}

// Pin returns true if the Cupdate pin label is set to true.
func (l Labels) Pin() bool {
	v, _ := l.oneOf("config.cupdate/pin", "cupdate.config.pin")
	return v == "true"
}

// StayOnCurrentMajor returns true if the Cupdate stay-on-current-major label is
// set to true.
func (l Labels) StayOnCurrentMajor() bool {
	v, _ := l.oneOf("config.cupdate/stay-on-current-major", "cupdate.config.stay-on-current-major")
	return v == "true"
}

// StayBelow returns the a semver if the Cupdate stay-below label is set.
// Returns an error if the value is set, but invalid.
func (l Labels) StayBelow() (*semver.Version, error) {
	v, ok := l.oneOf("config.cupdate/stay-below", "cupdate.config.stay-below")
	if !ok {
		return nil, nil
	}

	version, err := semver.ParseVersion(v)
	if err != nil {
		return nil, fmt.Errorf("label: %w", err)
	}

	if version.Prerelease != "" || version.Suffix != "" {
		return nil, fmt.Errorf("label: invalid semver - only specify release")
	}

	return version, nil
}

func (l Labels) oneOf(keys ...string) (string, bool) {
	if l == nil {
		return "", false
	}

	for _, k := range keys {
		if v, ok := l[k]; ok {
			return v, ok
		}
	}

	return "", false
}

// RemoveUnsupported removes unsupported labels.
func (l Labels) RemoveUnsupported() Labels {
	clone := maps.Clone(l)
	for k := range l {
		if !strings.HasPrefix(k, "config.cupdate/") && !strings.HasPrefix(k, "cupdate.config.") {
			delete(clone, k)
		}
	}
	return clone
}
