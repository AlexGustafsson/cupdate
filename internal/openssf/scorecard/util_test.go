package scorecard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryIsSupported(t *testing.T) {
	testCases := []struct {
		Repository string
		Expected   bool
	}{
		{
			Repository: "github.com/home-assistant/core",
			Expected:   true,
		},
		{
			Repository: "gitlab.com/baserow/baserow",
			Expected:   true,
		},
		{
			Repository: "https://github.com/home-assistant/core",
			Expected:   false,
		},
		{
			Repository: "http://gitlab.com/baserow/baserow",
			Expected:   false,
		},
		{
			Repository: "quay.io/some/project",
			Expected:   false,
		},
		{
			Repository: "192.168.1.248:8080/zot/project",
			Expected:   false,
		},
		{
			Repository: "",
			Expected:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Repository, func(t *testing.T) {
			assert.Equal(t, testCase.Expected, RepositoryIsSupported(testCase.Repository))
		})
	}
}
