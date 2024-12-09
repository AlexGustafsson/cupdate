package htmlutil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveReferences(t *testing.T) {
	input := `<div>
	<a href="/foo/bar">Hello</a>
	<img src="../../foo.png" />
	<img src="foo.png" />
	<img src="https://example.com/2.png" />
</div>`

	expected := `<div>
	<a href="https://example.com/foo/bar">Hello</a>
	<img src="https://example.com/foo.png"/>
	<img src="https://example.com/a/b/foo.png"/>
	<img src="https://example.com/2.png"/>
</div>`

	base, err := url.Parse("https://example.com/a/b/")
	require.NoError(t, err)

	actual, err := ResolveReferences(input, base)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
