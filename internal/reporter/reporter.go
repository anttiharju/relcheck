package reporter

import (
	"fmt"
	"strings"

	"github.com/anttiharju/relcheck/internal/color"
	"github.com/anttiharju/relcheck/internal/markdown/link"
)

// Reporter manages formatting and display of link check results
type Reporter struct {
	Colors     color.Palette
	Verbose    bool
	ErrorCount int
}

func New(colors color.Palette, verbose bool) *Reporter {
	return &Reporter{
		Colors:  colors,
		Verbose: verbose,
	}
}

// FileNotFound reports a file that doesn't exist
func (r *Reporter) FileNotFound(filename string) {
	fmt.Printf("%sError:%s %sFile not found: %s%s\n",
		r.Colors.Bold, r.Colors.Reset, r.Colors.Red, r.Colors.Reset, filename)

	r.ErrorCount++
}

// ReportProcessingError reports an error processing a file
func (r *Reporter) ReportProcessingError(filename string, err error) {
	fmt.Printf("%sError:%s Could not process file %s: %v\n",
		r.Colors.Bold, r.Colors.Reset, filename, err)

	r.ErrorCount++
}

// ReportNoLinks reports when no links were found (verbose only)
func (r *Reporter) ReportNoLinks(filename string) {
	if r.Verbose {
		fmt.Printf("%s✓%s %s: %sno relative links%s\n",
			r.Colors.Green, r.Colors.Reset, filename, r.Colors.Gray, r.Colors.Reset)
	}
}

// ReportBrokenLink reports a broken link
func (r *Reporter) ReportBrokenLink(filename string, brokenLink link.Link, errorType string, lineContent string) {
	fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (%s):%s\n",
		r.Colors.Bold, filename, brokenLink.Line, brokenLink.Column,
		r.Colors.Reset, r.Colors.Red, errorType, r.Colors.Reset)

	// Show the line content
	fmt.Println(lineContent)

	// Show the pointer to the link
	fmt.Printf("%s%s%s\n", r.Colors.Yellow, strings.Repeat(" ", brokenLink.Column-1)+"^", r.Colors.Reset)

	r.ErrorCount++
}

// ReportValidLinks reports the number of valid links found
func (r *Reporter) ReportValidLinks(filename string, count int, hasBrokenLinks bool) {
	if !r.Verbose || count == 0 {
		return
	}

	// Format the count text based on singular/plural
	countText := "1 valid relative link"
	if count > 1 {
		countText = fmt.Sprintf("%d valid relative links", count)
	}

	// Format the output based on whether broken links were found
	if !hasBrokenLinks {
		fmt.Printf("%s✓%s %s: found %s\n",
			r.Colors.Green, r.Colors.Reset, filename, countText)
	} else {
		fmt.Printf("%s%s: also found %s%s\n",
			r.Colors.Gray, filename, countText, r.Colors.Reset)
	}
}

// ReportSuccess reports total success when all links are valid
func (r *Reporter) ReportSuccess() {
	if r.Verbose && r.ErrorCount == 0 {
		fmt.Printf("%s✓%s %sAll relative links are valid!%s\n",
			r.Colors.Green, r.Colors.Reset, r.Colors.Bold, r.Colors.Reset)
	}
}
