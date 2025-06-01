package scan

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anttiharju/relcheck/internal/markdown/anchor"
	"github.com/anttiharju/relcheck/internal/markdown/link"
)

// Result stores both links and anchors from a single file scan
type Result struct {
	Links   []link.Link
	Anchors []string
}

//nolint:gochecknoglobals
var scanCache = make(map[string]Result)

//nolint:cyclop,funlen
func File(filename string) (Result, error) {
	// Check if we've already scanned this file
	if result, ok := scanCache[filename]; ok {
		return result, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return Result{}, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var links []link.Link

	var anchors []string

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	lineNumber := 0
	anchorCount := make(map[string]int)

	// Markdown link regex pattern - matches relative links with ./ or ../ prefixes
	relativeLinkPattern := regexp.MustCompile(`\]\(\.[^)]*\)`)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Check for code block
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock

			continue
		}

		// Skip content in code blocks
		if inCodeBlock {
			continue
		}

		// 1. Extract links from this line
		matches := relativeLinkPattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			start, end := match[0], match[1]
			// Extract URL without ]( and )
			urlText := line[start+2 : end-1]
			// Column position is start+2 (matching bash script)
			colPosition := start + 2

			path, anchorText := link.SplitLinkAndAnchor(urlText)

			links = append(links, link.Link{
				URL:    urlText,
				Line:   lineNumber,
				Column: colPosition + 1, // +1 because columns start at 1
				Path:   path,
				Anchor: anchorText,
			})
		}

		// 2. Look for headings in this line
		if strings.HasPrefix(line, "#") {
			// Match at least one # followed by a space
			if !regexp.MustCompile(`^#{1,6} `).MatchString(line) {
				continue
			}

			// Extract heading text without the leading #s
			heading := regexp.MustCompile(`^#+[ \t]+`).ReplaceAllString(line, "")
			// Remove trailing spaces
			heading = strings.TrimRight(heading, " \t")

			// Generate anchor
			anchorText := anchor.GenerateAnchor(heading)

			// Handle duplicate anchors
			if count := anchorCount[anchorText]; count > 0 {
				anchors = append(anchors, fmt.Sprintf("%s-%d", anchorText, count))
			} else {
				anchors = append(anchors, anchorText)
			}

			// Increment the counter for this anchor
			anchorCount[anchorText]++
		}
	}

	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("error scanning file: %w", err)
	}

	result := Result{
		Links:   links,
		Anchors: anchors,
	}

	// Cache the result
	scanCache[filename] = result

	return result, nil
}
