package htmlutil

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Resolve hrefs and srcs in input.
func ResolveReferences(input string, base *url.URL) (string, error) {
	node, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	err = Walk(node, func(node *html.Node) error {
		href, ok := GetAttr(node, "href")
		if ok {
			SetAttr(node, "href", resolveReference(href, base))
		}

		src, ok := GetAttr(node, "src")
		if ok {
			SetAttr(node, "src", resolveReference(src, base))
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := html.Render(&buffer, node); err != nil {
		return "", err
	}

	return strings.TrimSuffix(strings.TrimPrefix(buffer.String(), "<html><head></head><body>"), "</body></html>"), nil
}

func resolveReference(ref string, base *url.URL) string {
	u, err := base.Parse(ref)
	if err != nil {
		fmt.Println(ref, err)
		return ref
	}

	return u.String()
}