package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageReference(t *testing.T) {
	testCases := []struct {
		SpecImage     string
		StatusImage   string
		StatusImageID string
		Expected      string
		ExpectErr     bool
	}{
		{
			SpecImage: "mongo:4",
			// No status present - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:   "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			StatusImage: "mongo",
			// Status image present, but not as detailed - use image from spec
			Expected: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo",
			// Status image id present, but not as detailed - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, but without tag - use tag from spec
			// and sha from status
			Expected: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, with tag - use image from status
			Expected: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo",
			StatusImage:   "mongo:latest",
			StatusImageID: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, with tag - use image from status
			Expected: "mongo:latest@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "aW52YWxpZA==",
			// Status image id present, but not a valid reference - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "aW52YWxpZA==",
			StatusImage:   "aW52YWxpZA==",
			// Status image and image id present, but not a valid reference - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage: "",
			// Image not present in spec
			ExpectErr: true,
		},
		{
			SpecImage: "aW52YWxpZA==",
			// Image present in spec, but invalid
			ExpectErr: true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := getImageReference(testCase.SpecImage, testCase.StatusImage, testCase.StatusImageID)
			assert.Equal(t, testCase.Expected, actual.String())
			if testCase.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
