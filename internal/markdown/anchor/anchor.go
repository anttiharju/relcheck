package anchor

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Anchor represents a heading anchor in a markdown file
type Anchor struct {
	Name  string
	Line  int
	Level int // Heading level (1-6)
}

// Extract finds all heading anchors in a markdown file
func Extract(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var anchors []string

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	anchorCount := make(map[string]int)

	for scanner.Scan() {
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

		// Match headings
		if strings.HasPrefix(line, "#") {
			// Match at least one # followed by a space
			if !regexp.MustCompile(`^#{1,6} `).MatchString(line) {
				continue
			}

			// Extract heading text without the leading #s
			heading := regexp.MustCompile(`^#+[ \t]+`).ReplaceAllString(line, "")
			// Remove trailing spaces
			heading = strings.TrimRight(heading, " \t")

			anchor := GenerateAnchor(heading)

			// Handle duplicate anchors
			if count := anchorCount[anchor]; count > 0 {
				anchors = append(anchors, fmt.Sprintf("%s-%d", anchor, count))
			} else {
				anchors = append(anchors, anchor)
			}

			// Increment the counter for this anchor
			anchorCount[anchor]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return anchors, nil
}

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
