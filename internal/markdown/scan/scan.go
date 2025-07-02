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

type Result struct {
	Links     []link.Link
	Anchors   []string
	LineCount int
}

//nolint:gochecknoglobals
var scanCache = make(map[string]Result)

var (
	relativeLinkPattern = regexp.MustCompile(`\]\(\.[^)"']*(?:"[^"]*"|'[^']*')?\)`)
	headingPattern      = regexp.MustCompile(`^#{1,6} `)
	headingTextPattern  = regexp.MustCompile(`^#+[ \t]+`)
)

func File(filepath string) (Result, error) {
	// Check cache first
	if result, ok := scanCache[filepath]; ok {
		return result, nil
	}

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	// Scan the file
	result, err := scanFile(file)
	if err != nil {
		return Result{}, err
	}

	// Cache the result
	scanCache[filepath] = result

	return result, nil
}

func scanFile(file *os.File) (Result, error) {
	links := []link.Link{}
	anchors := []string{}

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	lineNumber := 0
	anchorCount := make(map[string]int)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if hasCodeBlockMarker(line) {
			inCodeBlock = !inCodeBlock

			continue
		}

		if inCodeBlock {
			continue
		}

		if extractHeading(&anchors, line, anchorCount) {
			continue // skip link extraction if line is a heading
		}

		extractLink(&links, line, lineNumber)
	}

	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("error scanning file: %w", err)
	}

	return Result{
		Links:     links,
		Anchors:   anchors,
		LineCount: lineNumber,
	}, nil
}

func hasCodeBlockMarker(line string) bool {
	trimmedLine := strings.TrimLeft(line, " \t")

	return strings.HasPrefix(trimmedLine, "```")
}

func extractLink(links *[]link.Link, line string, lineNumber int) {
	matches := relativeLinkPattern.FindAllStringIndex(line, -1)
	for _, match := range matches {
		start, end := match[0], match[1]
		// Extract URL without ]( and )
		rawURL := line[start+2 : end-1]
		// Skip ](
		colPosition := start + 2

		// Handle alt text in quotes if present
		urlText := rawURL
		if idx := strings.IndexAny(rawURL, "\"'"); idx != -1 {
			// Only take the part before the quote
			urlText = strings.TrimSpace(rawURL[:idx])
		}

		path, anchorText := link.SplitLinkAndAnchor(urlText)

		*links = append(*links, link.Link{
			URL:         urlText,
			Line:        lineNumber,
			Column:      colPosition + 1, // +1 because columns start at 1
			Path:        path,
			Anchor:      anchorText,
			LineContent: line,
		})
	}
}

func extractHeading(anchors *[]string, line string, anchorCount map[string]int) bool {
	if !strings.HasPrefix(line, "#") {
		return false
	}

	// Match 1 to 6 #s followed by a space
	if !headingPattern.MatchString(line) {
		return false
	}

	// Extract heading text without the leading #s
	heading := headingTextPattern.ReplaceAllString(line, "")
	// Remove trailing spaces
	heading = strings.TrimRight(heading, " \t")

	// Generate anchor
	anchorText := anchor.GenerateAnchor(heading)

	// Handle duplicate anchors
	if count := anchorCount[anchorText]; count > 0 {
		*anchors = append(*anchors, fmt.Sprintf("%s-%d", anchorText, count))
	} else {
		*anchors = append(*anchors, anchorText)
	}

	// Increment the counter for this anchor
	anchorCount[anchorText]++

	return true
}
