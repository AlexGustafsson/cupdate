package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	host, owner, repository, path, ok := ParseURL("https://github.com/docker-library/mongo.git#39c6083702fb2ff810e7a6a916b1eadf54825acd:6.0")
	require.True(t, ok)

	assert.Equal(t, "https://github.com", host)
	assert.Equal(t, "docker-library", owner)
	assert.Equal(t, "mongo", repository)
	assert.Equal(t, "/", path)
}
