package github

import (
	"net/url"
	"regexp"
	"strings"
)

func ParseURL(u string) (string, string, string, string, bool) {
	url, err := url.Parse(u)
	if err != nil {
		return "", "", "", "", false
	}

	parts := strings.Split(strings.TrimPrefix(url.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", "", "", false
	}

	path := "/"
	if len(parts) > 2 {
		path = strings.Join(parts[2:], "/")
	}

	repository := parts[1]
	if strings.Contains(repository, ".git") {
		repository = regexp.MustCompile(`\.git.*$`).ReplaceAllString(repository, "")
	}

	return url.Scheme + "://" + url.Host, parts[0], repository, path, true
}
