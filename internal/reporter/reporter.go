package reporter

import (
	"fmt"
	"strings"

	"github.com/anttiharju/relcheck/internal/color"
	"github.com/anttiharju/relcheck/internal/markdown/link"
)

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

func (r *Reporter) FileNotFound(filename string) {
	fmt.Printf("%sError:%s %sFile not found: %s%s\n",
		r.Colors.Bold, r.Colors.Reset, r.Colors.Red, r.Colors.Reset, filename)

	r.ErrorCount++
}

func (r *Reporter) ScanError(filename string, err error) {
	fmt.Printf("%sError:%s Could not process file %s: %v\n",
		r.Colors.Bold, r.Colors.Reset, filename, err)

	r.ErrorCount++
}

func (r *Reporter) NoLinks(filename string) {
	if r.Verbose {
		fmt.Printf("%s✓%s %s: %sno relative links%s\n",
			r.Colors.Green, r.Colors.Reset, filename, r.Colors.Gray, r.Colors.Reset)
	}
}

func (r *Reporter) BrokenLink(filename string, brokenLink link.Link, errorType string, lineContent string) {
	fmt.Printf("%s%s:%d:%d:%s %sbroken relative link (%s):%s\n",
		r.Colors.Bold, filename, brokenLink.Line, brokenLink.Column,
		r.Colors.Reset, r.Colors.Red, errorType, r.Colors.Reset)

	fmt.Println(lineContent)
	fmt.Printf("%s%s%s\n", r.Colors.Yellow, strings.Repeat(" ", brokenLink.Column-1)+"^", r.Colors.Reset)

	r.ErrorCount++
}

func (r *Reporter) ValidLinks(filename string, count int, hasBrokenLinks bool) {
	if !r.Verbose || count == 0 {
		return
	}

	// Singular/plural formatting
	countText := "1 valid relative link"
	if count > 1 {
		countText = fmt.Sprintf("%d valid relative links", count)
	}

	// Print with "also has" if there were broken links
	if !hasBrokenLinks {
		fmt.Printf("%s✓%s %s: %s\n",
			r.Colors.Green, r.Colors.Reset, filename, countText)
	} else {
		fmt.Printf("%s%s: also has %s%s\n",
			r.Colors.Gray, filename, countText, r.Colors.Reset)
	}
}

func (r *Reporter) Success() {
	if r.Verbose && r.ErrorCount == 0 {
		fmt.Printf("%s✓%s %sAll relative links are valid!%s\n",
			r.Colors.Green, r.Colors.Reset, r.Colors.Bold, r.Colors.Reset)
	}
}
