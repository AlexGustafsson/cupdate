package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameFromAPI(t *testing.T) {
	testCases := []struct {
		Path     string
		Expected string
	}{
		{
			Path:     "/v2/arm-research/smarter/smarter-device-manager/tags/list",
			Expected: "arm-research/smarter/smarter-device-manager",
		},
		{
			Path:     "/v2/alexgustafsson/cupdate/vulndb/tags/list",
			Expected: "alexgustafsson/cupdate/vulndb",
		},
		{
			// Blobs OK in name
			Path:     "/v2/test/blobs/blobs/ref",
			Expected: "test/blobs",
		},
		{
			// Empty
			Path:     "",
			Expected: "",
		},
		{
			// Too few path segments
			Path:     "/",
			Expected: "",
		},
		{
			// Too few path segments
			Path:     "/v2",
			Expected: "",
		},
		{
			// Invalid path
			Path:     "v2/alexgustafsson/cupdate/vulndb/tags/list",
			Expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Path, func(t *testing.T) {
			actual := NameFromAPI(testCase.Path)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}
