package check

import (
	"net/url"

	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/fileutils"
	"github.com/anttiharju/relcheck/internal/markdown/anchor"
	"github.com/anttiharju/relcheck/internal/markdown/link"
	"github.com/anttiharju/relcheck/internal/markdown/scan"
	"github.com/anttiharju/relcheck/internal/reporter"
)

// RelativeLinksAndAnchors checks markdown files for broken relative links
func RelativeLinksAndAnchors(verbose, forceColors bool, files []string) exitcode.Exitcode {
	report := reporter.New(verbose, forceColors)
	exitCode := exitcode.Success

	for _, filepath := range files {
		fileExitCode := processFile(filepath, report)
		if fileExitCode != exitcode.Success {
			exitCode = fileExitCode
		}
	}

	// Report success if all links are valid
	if exitCode == exitcode.Success {
		report.Success()
	}

	return exitCode
}

// processFile handles a single markdown file
func processFile(filepath string, report *reporter.Reporter) exitcode.Exitcode {
	// Check if the file exists
	if !fileutils.FileExists(filepath) {
		report.FileNotFound(filepath)

		return exitcode.BrokenLinks
	}

	// Scan the file for links
	scanResult, err := scan.File(filepath)
	if err != nil {
		report.ScanError(filepath, err)

		return exitcode.BrokenLinks
	}

	// If no links are found, report and continue
	if len(scanResult.Links) == 0 {
		report.NoLinks(filepath)

		return exitcode.Success
	}

	return processLinks(filepath, scanResult, report)
}

// processLinks validates all links in a file
func processLinks(filepath string, scanResult scan.Result, report *reporter.Reporter) exitcode.Exitcode {
	brokenLinksFound := false
	validLinksCount := 0

	for _, link := range scanResult.Links {
		valid := validateLink(filepath, link, report)
		if valid {
			validLinksCount++
		} else {
			brokenLinksFound = true
		}
	}

	// Report valid links if any
	report.ValidLinks(filepath, validLinksCount, brokenLinksFound)

	if brokenLinksFound {
		return exitcode.BrokenLinks
	}

	return exitcode.Success
}

// validateLink checks if a single link is valid
func validateLink(filepath string, link link.Link, report *reporter.Reporter) bool {
	// URL-decode the link path
	decodedPath, err := url.QueryUnescape(link.Path)
	if err != nil {
		report.ScanError(filepath, err)

		return false
	}

	// Construct the full path relative to the file's location
	fullpath := fileutils.ResolveRelativePath(filepath, decodedPath)

	// Check if target file exists
	if !fileutils.FileExists(fullpath) {
		lineContent, _ := fileutils.GetLineContent(filepath, link.Line)
		report.BrokenLink(filepath, link, "target not found", lineContent)

		return false
	}

	// If an anchor exists, check if it's valid
	if link.Anchor != "" {
		return validateAnchor(filepath, fullpath, link, report)
	}

	// If we got here, the link is valid
	return true
}

// validateAnchor checks if an anchor in a link is valid
func validateAnchor(filepath, fullpath string, link link.Link, report *reporter.Reporter) bool {
	// Get anchors from target file (uses cache when possible)
	targetScan, err := scan.File(fullpath)
	if err != nil {
		report.ScanError(filepath, err)

		return false
	}

	// Check if the anchor exists
	if !anchor.Contains(targetScan.Anchors, link.Anchor) {
		lineContent, _ := fileutils.GetLineContent(filepath, link.Line)
		report.BrokenLink(filepath, link, "anchor not found", lineContent)

		return false
	}

	return true
}
