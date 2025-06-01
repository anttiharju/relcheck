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
	Links   []link.Link
	Anchors []string
}

//nolint:gochecknoglobals
var scanCache = make(map[string]Result)

var (
	relativeLinkPattern = regexp.MustCompile(`\]\(\.[^)]*\)`) // catch ./ or ../ prefixes
	headingPattern      = regexp.MustCompile(`^#{1,6} `)
	headingTextPattern  = regexp.MustCompile(`^#+[ \t]+`)
)

func File(filepath string) (Result, error) {
	// Check cache first
	if result, ok := scanCache[filepath]; ok {
		return result, nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

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

		// Check for code block markers
		if isCodeBlockMarker(line) {
			inCodeBlock = !inCodeBlock

			continue
		}

		// Skip content in code blocks
		if inCodeBlock {
			continue
		}

		// Process links and headings
		extractLinks(&links, line, lineNumber)
		extractHeadings(&anchors, line, anchorCount)
	}

	if err := scanner.Err(); err != nil {
		return Result{}, fmt.Errorf("error scanning file: %w", err)
	}

	return Result{
		Links:   links,
		Anchors: anchors,
	}, nil
}

func isCodeBlockMarker(line string) bool {
	trimmedLine := strings.TrimLeft(line, " ")

	return strings.HasPrefix(trimmedLine, "```")
}

func extractLinks(links *[]link.Link, line string, lineNumber int) {
	matches := relativeLinkPattern.FindAllStringIndex(line, -1)
	for _, match := range matches {
		start, end := match[0], match[1]
		// Extract URL without ]( and )
		urlText := line[start+2 : end-1]
		// Column position is start+2
		colPosition := start + 2

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

func extractHeadings(anchors *[]string, line string, anchorCount map[string]int) {
	if !strings.HasPrefix(line, "#") {
		return
	}

	// Match at least one # followed by a space
	if !headingPattern.MatchString(line) {
		return
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
}
