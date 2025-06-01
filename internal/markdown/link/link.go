package link

import (
	"strings"
)

// Link represents a relative link in a markdown file
type Link struct {
	URL     string
	Line    int
	Column  int
	Path    string // Resolved path
	Anchor  string // Anchor part if present
	IsValid bool   // Validation status
}

// SplitLinkAndAnchor splits a link into path and anchor parts
func SplitLinkAndAnchor(link string) (string, string) {
	parts := strings.SplitN(link, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return link, ""
}
