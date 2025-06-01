package anchor

import (
	"regexp"
	"strings"
)

// GenerateAnchor creates a GitHub-compatible anchor from heading text
func GenerateAnchor(heading string) string {
	// Convert to lowercase
	anchor := strings.ToLower(heading)

	// Remove anything that's not alphanumeric, space, or hyphen
	anchor = regexp.MustCompile(`[^a-z0-9 -]`).ReplaceAllString(anchor, "")

	// Convert spaces to hyphens
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// Replace multiple hyphens with single hyphen
	anchor = regexp.MustCompile(`-+`).ReplaceAllString(anchor, "-")

	// Trim trailing hyphens
	anchor = strings.TrimRight(anchor, "-")

	return anchor
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
