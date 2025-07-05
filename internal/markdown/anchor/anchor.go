package anchor

import (
	"regexp"
	"slices"
	"strings"
)

var validCharactersPattern = regexp.MustCompile(`[^a-z0-9 _-]`)

var multipleHyphensPattern = regexp.MustCompile(`-+`)

func GenerateAnchor(heading string) string {
	anchor := strings.ToLower(heading)

	// Remove anything that's not alphanumeric, space, or hyphen
	anchor = validCharactersPattern.ReplaceAllString(anchor, "")

	// Convert spaces to hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Replace multiple hyphens with single hyphen
	anchor = multipleHyphensPattern.ReplaceAllString(anchor, "-")

	// Trim trailing hyphens
	anchor = strings.TrimRight(anchor, "-")

	return anchor
}

func Exists(target []string, source string) bool {
	return slices.Contains(target, source)
}
