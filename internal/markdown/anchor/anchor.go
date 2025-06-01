package anchor

import (
	"regexp"
	"slices"
	"strings"
)

// alphanumeric, space and hyphen
var alphanumericSpaceHyphenPattern = regexp.MustCompile(`[^a-z0-9 -]`)

var multipleHyphensPattern = regexp.MustCompile(`-+`)

// GenerateAnchor creates a GitHub-compatible anchor from heading text
func GenerateAnchor(heading string) string {
	// Convert to lowercase
	anchor := strings.ToLower(heading)

	// Remove anything that's not alphanumeric, space, or hyphen
	anchor = alphanumericSpaceHyphenPattern.ReplaceAllString(anchor, "")

	// Convert spaces to hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Replace multiple hyphens with single hyphen
	anchor = multipleHyphensPattern.ReplaceAllString(anchor, "-")

	// Trim trailing hyphens
	anchor = strings.TrimRight(anchor, "-")

	return anchor
}

func Exists(source string, target []string) bool {
	return slices.Contains(target, source)
}
