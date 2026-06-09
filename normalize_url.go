package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func normalizeURL(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", ErrInvalidURL
	}

	stripped := fmt.Sprintf("%s%s", parsed.Host, path.Clean(strings.ReplaceAll(parsed.Path, `\`, "/")))
	stripped = strings.ToLower(stripped)

	return stripped, nil
}
