package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Terminal colors
const (
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	gray   = "\033[90m"
	reset  = "\033[0m"
)

// CLI flags
type options struct {
	verbose    bool
	forceColor bool
	files      []string
}

func main() {
	opts := parseFlags()

	// If no files provided, show usage
	if len(opts.files) == 0 {
		fmt.Println("Usage: relcheck [--verbose] [--color=always] <file1.md> [file2.md] ...")
		fmt.Println("   or: relcheck [--verbose] [--color=always] run  (to check all *.md files in Git)")
		os.Exit(1)
	}

	// Determine terminal colors
	useColors := isTerminal() || opts.forceColor
	colors := getColorScheme(useColors)

	exitCode := 0
	for _, file := range opts.files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("%sError:%s %sFile not found: %s%s\n", colors.bold, colors.reset, colors.red, colors.reset, file)
			exitCode = 1
			continue
		}

		// Get directory of the current file to resolve relative paths
		dir := filepath.Dir(file)

		// Extract links from the markdown file
		links, err := extractRelativeLinks(file)
		if err != nil {
			fmt.Printf("%sError:%s Could not process file %s: %v\n", colors.bold, colors.reset, file, err)
			exitCode = 1
			continue
		}

		// If no links are found, continue to the next file
		if len(links) == 0 {
			if opts.verbose {
				fmt.Printf("%s✓%s %s: %sno relative links%s\n", colors.green, colors.reset, file, colors.gray, colors.reset)
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
				fmt.Printf("%sError:%s Could not decode URL %s: %v\n", colors.bold, colors.reset, linkPath, err)
				continue
			}

			// Construct the full path relative to the file's location
			fullPath := filepath.Join(dir, decodedLink)

			// Use base filename only for error messages (to match bash script output)
			baseFilename := filepath.Base(file)

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				// Print the file location in bold - match exact format from the bash script
				fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (file not found):%s\n",
					colors.bold, baseFilename, link.line, link.col, colors.reset, colors.red, colors.reset)

				// Extract the line content for context
				lineContent, _ := getLineContent(file, link.line)
				fmt.Println(lineContent)

				// Print line content with yellow indicator pointing to the link position
				fmt.Printf("%s%s%s\n", colors.yellow, strings.Repeat(" ", link.col-1)+"^", colors.reset)
				brokenLinksFound = true
			} else if linkAnchor != "" {
				// If an anchor exists, check if it's valid
				anchors, err := getMarkdownAnchors(fullPath)
				if err != nil {
					fmt.Printf("%sError:%s Could not extract anchors from %s: %v\n", colors.bold, colors.reset, fullPath, err)
					continue
				}

				if !contains(anchors, linkAnchor) {
					fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (anchor not found):%s\n",
						colors.bold, baseFilename, link.line, link.col, colors.reset, colors.red, colors.reset)
					lineContent, _ := getLineContent(file, link.line)
					fmt.Println(lineContent)
					fmt.Printf("%s%s%s\n", colors.yellow, strings.Repeat(" ", link.col-1)+"^", colors.reset)
					brokenLinksFound = true
				} else {
					validLinksCount++
				}
			} else {
				validLinksCount++
			}
		}

		// If verbose mode and we have valid links, report them
		if opts.verbose && validLinksCount > 0 {
			if !brokenLinksFound {
				if validLinksCount == 1 {
					fmt.Printf("%s✓%s %s: found 1 valid relative link\n", colors.green, colors.reset, file)
				} else {
					fmt.Printf("%s✓%s %s: found %d valid relative links\n", colors.green, colors.reset, file, validLinksCount)
				}
			} else {
				if validLinksCount == 1 {
					fmt.Printf("%s%s: also found 1 valid relative link%s\n", colors.gray, file, colors.reset)
				} else {
					fmt.Printf("%s%s: also found %d valid relative links%s\n", colors.gray, file, validLinksCount, colors.reset)
				}
			}
		}

		if brokenLinksFound {
			exitCode = 1
		}
	}

	// Show success message if all links are valid, but only in verbose mode
	if exitCode == 0 && opts.verbose {
		fmt.Printf("%s✓%s %sAll relative links are valid!%s\n", colors.green, colors.reset, colors.bold, colors.reset)
	}

	os.Exit(exitCode)
}

// Link represents a relative link in a markdown file
type Link struct {
	url  string
	line int
	col  int
}

// Markdown link regex pattern - matches relative links like [text](./path) or [text](../path)
var relativeLinkPattern = regexp.MustCompile(`\]\((\.\.?/[^)]*)\)`)

// extracts relative links from a markdown file
func extractRelativeLinks(filename string) ([]Link, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
		matches := relativeLinkPattern.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			// match[0], match[1] is the entire match
			// match[2], match[3] is the capture group (\.\.?/[^)]*)
			start, end := match[2], match[3]
			// Extract URL from the capture group
			url := line[start:end]
			// Calculate column position - we need to adjust to where the link actually starts
			urlStartInMarkdown := strings.LastIndex(line[:start], "(") + 1
			links = append(links, Link{
				url:  url,
				line: lineNumber,
				col:  urlStartInMarkdown + 1, // +1 because columns start at 1
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
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
		return nil, err
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
			// Extract heading text without leading #s
			headingParts := strings.SplitN(line, " ", 2)
			if len(headingParts) < 2 {
				continue
			}

			heading := strings.TrimSpace(headingParts[1])

			// Convert to GitHub-style anchor
			anchor := strings.ToLower(heading)
			// Remove anything that's not alphanumeric, space, or hyphen
			re := regexp.MustCompile(`[^a-z0-9 -]`)
			anchor = re.ReplaceAllString(anchor, "")
			// Convert spaces to hyphens
			anchor = strings.ReplaceAll(anchor, " ", "-")
			// Replace multiple hyphens with single hyphen
			re = regexp.MustCompile(`-+`)
			anchor = re.ReplaceAllString(anchor, "-")
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
		return nil, err
	}

	return anchors, nil
}

// gets content of a specific line from a file
func getLineContent(filename string, lineNumber int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
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

// ColorScheme holds ANSI color codes
type ColorScheme struct {
	bold, red, yellow, green, gray, reset string
}

// returns color scheme based on terminal capabilities
func getColorScheme(useColors bool) ColorScheme {
	if useColors {
		return ColorScheme{bold, red, yellow, green, gray, reset}
	}
	return ColorScheme{"", "", "", "", "", ""}
}

// checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// parses command-line flags
func parseFlags() options {
	// Define the flags
	flagSet := flag.NewFlagSet("relcheck", flag.ExitOnError)
	verbose := flagSet.Bool("verbose", false, "Enable verbose output")
	color := flagSet.Bool("color", false, "Force color output")

	// Special handling for --color=always format
	colorAlways := false
	for i, arg := range os.Args {
		if arg == "--color=always" {
			// Remove this arg to avoid flagSet parsing error
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			colorAlways = true
			break
		}
	}

	// Parse the remaining flags
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	opts := options{
		verbose:    *verbose,
		forceColor: *color || colorAlways,
		files:      []string{},
	}

	// Process remaining arguments
	args := flagSet.Args()
	for i, arg := range args {
		if i == 0 && arg == "run" {
			// Use git ls-files to find all markdown files
			cmd := exec.Command("git", "ls-files", "*.md", "*.markdown")
			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err == nil {
				scanner := bufio.NewScanner(&out)
				for scanner.Scan() {
					opts.files = append(opts.files, scanner.Text())
				}
			}
		} else {
			opts.files = append(opts.files, arg)
		}
	}

	return opts
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
