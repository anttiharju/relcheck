package link

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
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

// Markdown link regex pattern - matches relative links with ./ or ../ prefixes
var relativeLinkPattern = regexp.MustCompile(`\]\(\.[^)]*\)`)

// Extract finds all relative links in a markdown file
func Extract(filename string) ([]Link, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var links []Link

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check for code block
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock

			continue
		}

		// Skip lines in code blocks
		if inCodeBlock {
			continue
		}

		// Find relative links in the line
		matches := relativeLinkPattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			start, end := match[0], match[1]
			// Extract URL without ]( and )
			urlText := line[start+2 : end-1]
			// Column position is start+2 (to match bash script)
			colPosition := start + 2

			path, anchor := SplitLinkAndAnchor(urlText)

			links = append(links, Link{
				URL:    urlText,
				Line:   lineNumber,
				Column: colPosition + 1, // +1 because columns start at 1
				Path:   path,
				Anchor: anchor,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return links, nil
}

// SplitLinkAndAnchor splits a link into path and anchor parts
func SplitLinkAndAnchor(link string) (string, string) {
	parts := strings.SplitN(link, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return link, ""
}

// DecodeURL safely decodes a URL-encoded path
func DecodeURL(path string) (string, error) {
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		return "", fmt.Errorf("failed to decode URL path: %w", err)
	}

	return decodedPath, nil
}
