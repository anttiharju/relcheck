package anchor

import (
	"regexp"
	"slices"
	"strings"
)

var validCharactersPattern = regexp.MustCompile(`[^a-z0-9 _-]`)

var multipleHyphensPattern = regexp.MustCompile(`-+`)

func normalise(anchor string) string {
	normalised := strings.ToLower(anchor)

	// Remove anything that's not alphanumeric, space, underscore, or hyphen
	normalised = validCharactersPattern.ReplaceAllString(normalised, "")

	// Convert spaces to hyphens
	normalised = strings.ReplaceAll(normalised, " ", "-")

	// Replace multiple hyphens with single hyphen
	normalised = multipleHyphensPattern.ReplaceAllString(normalised, "-")

	// Trim trailing hyphens
	normalised = strings.TrimRight(normalised, "-")

	return normalised
}

func GenerateAnchor(heading string) string {
	return normalise(heading)
}

func Exists(target []string, source string) bool {
	normalizedSource := normalise(source)

	return slices.Contains(target, normalizedSource)
}
