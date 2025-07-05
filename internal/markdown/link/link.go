package link

import (
	"strings"
)

type Link struct {
	URL         string
	Line        int
	Column      int
	Path        string // Resolved path
	Anchor      string // Anchor part if present
	IsValid     bool
	LineContent string
}

func SplitLinkAndAnchor(link string) (string, string) {
	parts := strings.SplitN(link, "#", 2)
	if len(parts) == 2 {
		pathPart := parts[0]
		anchorPart := parts[1]

		// If the path part is empty, treat the link as relative to current directory
		if pathPart == "" {
			return ".", anchorPart
		}

		return pathPart, anchorPart
	}

	return link, ""
}
