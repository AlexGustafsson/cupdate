package htmlutil

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestGetAttr(t *testing.T) {
	testCases := []struct {
		HTML          string
		Key           string
		ExpectedValue string
		ExpectedBool  bool
	}{
		{
			HTML:          `<a href="https://example.com"><p>Link</p></a>`,
			Key:           "href",
			ExpectedValue: "https://example.com",
			ExpectedBool:  true,
		},
		{
			HTML:          `<a href=""><p>Link</p></a>`,
			Key:           "href",
			ExpectedValue: "",
			ExpectedBool:  true,
		},
		{
			HTML:          `<a href="foo" href="bar"><p>Link</p></a>`,
			Key:           "href",
			ExpectedValue: "foo",
			ExpectedBool:  true,
		},
		{
			HTML:          `<a><p>Link</p></a>`,
			Key:           "href",
			ExpectedValue: "",
			ExpectedBool:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.HTML, func(t *testing.T) {
			nodes, err := html.ParseFragment(strings.NewReader(testCase.HTML), nil)
			require.NoError(t, err)

			// HTML -> Head - Body -> Node
			node := nodes[0].FirstChild.NextSibling.FirstChild

			value, ok := GetAttr(node, testCase.Key)
			assert.Equal(t, testCase.ExpectedValue, value)
			assert.Equal(t, testCase.ExpectedBool, ok)
		})
	}
}

func TestSetAttr(t *testing.T) {
	testCases := []struct {
		HTML         string
		Key          string
		Value        string
		ExpectedHTML string
	}{
		{
			HTML:         `<a href="https://example.com"><p>Link</p></a>`,
			Key:          "href",
			Value:        "foo",
			ExpectedHTML: `<a href="foo"><p>Link</p></a>`,
		},
		{
			HTML:         `<a><p>Link</p></a>`,
			Key:          "href",
			Value:        "foo",
			ExpectedHTML: `<a href="foo"><p>Link</p></a>`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.HTML, func(t *testing.T) {
			nodes, err := html.ParseFragment(strings.NewReader(testCase.HTML), nil)
			require.NoError(t, err)

			// HTML -> Head - Body -> Node
			node := nodes[0].FirstChild.NextSibling.FirstChild

			SetAttr(node, testCase.Key, testCase.Value)

			var buffer bytes.Buffer
			html.Render(&buffer, node)
			assert.Equal(t, testCase.ExpectedHTML, buffer.String())
		})
	}
}
