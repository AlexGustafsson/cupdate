package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageReference(t *testing.T) {
	testCases := []struct {
		Image       string
		RepoDigests []string
		Expected    string
		ExpectErr   bool
	}{
		{
			Image: "mongo:4",
			// No digests or tags present - use image from container
			Expected: "mongo:4",
		},
		{
			Image: "mongo:4",
			RepoDigests: []string{
				"mongo:4.1@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
			},
			// Single digest available with tag - overwrite it
			Expected: "mongo:4@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
		},
		{
			Image: "mongo:4",
			RepoDigests: []string{
				"mongo@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
			},
			// Single digest available without tag - use it and the container's tag
			Expected: "mongo:4@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
		},
		{
			Image: "mongo:4",
			RepoDigests: []string{
				"ghcr.io/mongo/mongo@sha256:115a5eec6d9391912cd5d0b750b6b3f3886c2984e1ca5d51c4d9f430dc3c7b2e",
				"mongo@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
			},
			// Multiple digests available, without tags, use the one matching the image
			Expected: "mongo:4@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
		},
		{
			Image: "mongo@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
			RepoDigests: []string{
				"mongo@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
			},
			// No tag, just digest
			Expected: "mongo@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := getImageReference(testCase.Image, testCase.RepoDigests)
			assert.Equal(t, testCase.Expected, actual.String())
			if testCase.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
