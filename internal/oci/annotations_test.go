package oci

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	annotations := Annotations{
		"org.opencontainers.image.created": "2024-12-20 10:52:57+00:00",
	}

	expected := time.Date(2024, 12, 20, 10, 52, 57, 0, time.FixedZone("", 0))

	assert.True(t, annotations.Created().Equal(expected))
}
