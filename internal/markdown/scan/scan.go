package scan

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anttiharju/relcheck/internal/fileutils"
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
	headingAltPattern   = regexp.MustCompile(`^(-+|=+)\s*$`)
	markdownLinkPattern = regexp.MustCompile(`\[(.*?)\]\([^\)]*\)`)
)

func File(filepath string) (Result, error) {
	// Check cache first
	if result, ok := scanCache[filepath]; ok {
		return result, nil
	}

	// Check if path is a directory
	isDir, err := fileutils.IsDirectory(filepath)
	if err != nil {
		return Result{}, fmt.Errorf("failed to check path %s: %w", filepath, err)
	}

	if isDir {
		return Result{}, fmt.Errorf("path is a directory, not a file: %s", filepath)
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

func isHTMLCommentLine(line string, inHTMLComment bool) (bool, bool) {
	trimmed := strings.TrimSpace(line)

	// Entering a multi-line comment
	if strings.HasPrefix(trimmed, "<!--") && !strings.Contains(trimmed, "-->") {
		return true, true
	}
	// Inside a multi-line comment
	if inHTMLComment {
		if strings.Contains(trimmed, "-->") {
			return true, false
		}

		return true, true
	}
	// Single-line comment
	if strings.HasPrefix(trimmed, "<!--") && strings.Contains(trimmed, "-->") {
		return true, false
	}
	// Not a comment
	return false, inHTMLComment
}

//nolint:funlen // the function is pretty simple even if it is long
func scanFile(file *os.File) (Result, error) {
	links := []link.Link{}
	anchors := []string{}

	scanner := bufio.NewScanner(file)
	inCodeBlock := false
	inHTMLComment := false
	lineNumber := 0
	anchorCount := make(map[string]int)

	var previousLine string // Store the previous line for altHeading detection

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		skip, newInHTMLComment := isHTMLCommentLine(line, inHTMLComment)
		inHTMLComment = newInHTMLComment

		if skip {
			previousLine = line

			continue
		}

		if hasCodeBlockMarker(line) {
			inCodeBlock = !inCodeBlock
			previousLine = line

			continue
		}

		if inCodeBlock {
			previousLine = line

			continue
		}

		if extractAltHeading(&anchors, line, previousLine, anchorCount) {
			previousLine = line

			continue
		}

		if extractHeading(&anchors, line, anchorCount) {
			previousLine = line

			continue
		}

		extractLink(&links, line, lineNumber)
		previousLine = line
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

	// Must start with ```
	if !strings.HasPrefix(trimmedLine, "```") {
		return false
	}

	afterMarker := trimmedLine[3:]

	// Ensure code block is not closed on the same line
	if !strings.Contains(afterMarker, "```") {
		return true
	}

	return false
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

	// Remove markdown link syntax from heading
	heading = markdownLinkPattern.ReplaceAllString(heading, "$1")

	anchorText := anchor.GenerateAnchor(heading)

	// Handle duplicate anchors
	if count := anchorCount[anchorText]; count > 0 {
		*anchors = append(*anchors, fmt.Sprintf("%s-%d", anchorText, count))
	} else {
		*anchors = append(*anchors, anchorText)
	}

	anchorCount[anchorText]++

	return true
}

func extractAltHeading(anchors *[]string, currentLine string, previousLine string, anchorCount map[string]int) bool {
	if previousLine != "" && !strings.HasPrefix(currentLine, "#") && headingAltPattern.MatchString(currentLine) {
		// Remove trailing spaces
		heading := strings.TrimRight(previousLine, " \t")

		anchorText := anchor.GenerateAnchor(heading)

		// Handle duplicate anchors
		if count := anchorCount[anchorText]; count > 0 {
			*anchors = append(*anchors, fmt.Sprintf("%s-%d", anchorText, count))
		} else {
			*anchors = append(*anchors, anchorText)
		}

		anchorCount[anchorText]++

		return true
	}

	return false
}
