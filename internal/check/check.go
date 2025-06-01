package check

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/anttiharju/relcheck/internal/color"
	"github.com/anttiharju/relcheck/internal/exitcode"
)

// Link represents a relative link in a markdown file
type Link struct {
	url  string
	line int
	col  int
}

//nolint:gocognit,cyclop,funlen
func RelativeLinksAndAnchors(verbose, forceColors bool, files []string) exitcode.Exitcode {
	// Determine terminal colors
	colors := color.GetPalette(forceColors)

	exitCode := exitcode.Success

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("%sError:%s %sFile not found: %s%s\n", colors.Bold, colors.Reset, colors.Red, colors.Reset, file)

			exitCode = 1

			continue
		}

		// Get the directory of the current file to resolve relative paths
		dir := filepath.Dir(file)

		// Extract links from the markdown file
		links, err := extractRelativeLinks(file)
		if err != nil {
			fmt.Printf("%sError:%s Could not process file %s: %v\n", colors.Bold, colors.Reset, file, err)

			exitCode = 1

			continue
		}

		// If no links are found, continue to the next file
		if len(links) == 0 {
			if verbose {
				fmt.Printf("%s✓%s %s: %sno relative links%s\n", colors.Green, colors.Reset, file, colors.Gray, colors.Reset)
			}

			continue
		}

		brokenLinksFound := false
		validLinksCount := 0

		// Process each link
		for _, link := range links {
			linkPath, linkAnchor := splitLinkAndAnchor(link.url)

			// URL-decode the link path
			decodedLink, err := url.QueryUnescape(linkPath)
			if err != nil {
				fmt.Printf("%sError:%s Could not decode URL %s: %v\n", colors.Bold, colors.Reset, linkPath, err)

				continue
			}

			// Construct the full path relative to the file's location
			fullPath := filepath.Join(dir, decodedLink)

			// Use full file path for error messages to match bash script output
			//nolint:nestif
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				// Print the file location in bold - match exact format from the bash script
				fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (file not found):%s\n",
					colors.Bold, file, link.line, link.col, colors.Reset, colors.Red, colors.Reset)

				// Extract the line content for context
				lineContent, _ := getLineContent(file, link.line)
				fmt.Println(lineContent)

				// Print line content with yellow indicator pointing to the link position
				fmt.Printf("%s%s%s\n", colors.Yellow, strings.Repeat(" ", link.col-1)+"^", colors.Reset)

				brokenLinksFound = true
			} else if linkAnchor != "" {
				// If an anchor exists, check if it's valid
				anchors, err := getMarkdownAnchors(fullPath)
				if err != nil {
					fmt.Printf("%sError:%s Could not extract anchors from %s: %v\n", colors.Bold, colors.Reset, fullPath, err)

					continue
				}

				if !contains(anchors, linkAnchor) {
					fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (anchor not found):%s\n",
						colors.Bold, file, link.line, link.col, colors.Reset, colors.Red, colors.Reset)

					lineContent, _ := getLineContent(file, link.line)
					fmt.Println(lineContent)
					fmt.Printf("%s%s%s\n", colors.Yellow, strings.Repeat(" ", link.col-1)+"^", colors.Reset)

					brokenLinksFound = true
				} else {
					validLinksCount++
				}
			} else {
				validLinksCount++
			}
		}

		// If verbose mode and we have valid links, report them
		//nolint:nestif
		if verbose && validLinksCount > 0 {
			if !brokenLinksFound {
				if validLinksCount == 1 {
					fmt.Printf("%s✓%s %s: found 1 valid relative link\n", colors.Green, colors.Reset, file)
				} else {
					fmt.Printf("%s✓%s %s: found %d valid relative links\n", colors.Green, colors.Reset, file, validLinksCount)
				}
			} else {
				if validLinksCount == 1 {
					fmt.Printf("%s%s: also found 1 valid relative link%s\n", colors.Gray, file, colors.Reset)
				} else {
					fmt.Printf("%s%s: also found %d valid relative links%s\n", colors.Gray, file, validLinksCount, colors.Reset)
				}
			}
		}

		if brokenLinksFound {
			exitCode = 1
		}
	}

	// Show success message if all links are valid, but only in verbose mode
	if exitCode == exitcode.Success && verbose {
		fmt.Printf("%s✓%s %sAll relative links are valid!%s\n", colors.Green, colors.Reset, colors.Bold, colors.Reset)
	}

	return exitCode
}

// Markdown link regex pattern - matches relative links with ./ or ../ prefixes
var relativeLinkPattern = regexp.MustCompile(`\]\(\.[^)]*\)`)

// extracts relative links from a markdown file
func extractRelativeLinks(filename string) ([]Link, error) {
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
		// This uses the exact same regex pattern as the bash script
		matches := relativeLinkPattern.FindAllStringIndex(line, -1)
		for _, match := range matches {
			start, end := match[0], match[1]
			// Extract URL without ]( and )
			url := line[start+2 : end-1]
			// Column position is start+2 (to match the bash script)
			colPosition := start + 2
			links = append(links, Link{
				url:  url,
				line: lineNumber,
				col:  colPosition + 1, // +1 because columns start at 1
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return links, nil
}

// splits a link into path and anchor parts
func splitLinkAndAnchor(link string) (string, string) {
	parts := strings.SplitN(link, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return link, ""
}

// gets markdown anchors from a file
func getMarkdownAnchors(filename string) ([]string, error) {
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

		// Match headings with regex pattern similar to bash script
		if strings.HasPrefix(line, "#") {
			// Match at least one # followed by a space
			if !regexp.MustCompile(`^#{1,6} `).MatchString(line) {
				continue
			}

			// Extract heading text without the leading #s
			heading := regexp.MustCompile(`^#+[ \t]+`).ReplaceAllString(line, "")
			// Remove trailing spaces
			heading = strings.TrimRight(heading, " \t")

			// Convert to GitHub-style anchor
			anchor := strings.ToLower(heading)
			// Remove anything that's not alphanumeric, space, or hyphen
			anchor = regexp.MustCompile(`[^a-z0-9 -]`).ReplaceAllString(anchor, "")
			// Convert spaces to hyphens
			anchor = strings.ReplaceAll(anchor, " ", "-")
			// Replace multiple hyphens with single hyphen
			anchor = regexp.MustCompile(`-+`).ReplaceAllString(anchor, "-")
			// Trim trailing hyphens
			anchor = strings.TrimRight(anchor, "-")

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

// gets content of a specific line from a file
func getLineContent(filename string, lineNumber int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine == lineNumber {
			return scanner.Text(), nil
		}
	}

	return "", fmt.Errorf("line %d not found", lineNumber)
}

// checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
