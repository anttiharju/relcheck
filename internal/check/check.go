package check

import (
	"net/url"

	"github.com/anttiharju/relcheck/internal/color"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/fileutils"
	"github.com/anttiharju/relcheck/internal/markdown/anchor"
	"github.com/anttiharju/relcheck/internal/markdown/scanner"
	"github.com/anttiharju/relcheck/internal/reporter"
)

//nolint:cyclop,funlen
func RelativeLinksAndAnchors(verbose, forceColors bool, files []string) exitcode.Exitcode {
	// Initialize reporter with terminal colors
	colors := color.GetPalette(forceColors)
	reporter := reporter.New(colors, verbose)

	// Default exit code is success
	exitCode := exitcode.Success

	for _, file := range files {
		// Check if the file exists
		if !fileutils.FileExists(file) {
			reporter.ReportFileNotFound(file)

			exitCode = exitcode.BrokenLinks

			continue
		}

		// Scan the file for links
		scanResult, err := scanner.ScanFile(file)
		if err != nil {
			reporter.ReportProcessingError(file, err)

			exitCode = exitcode.BrokenLinks

			continue
		}

		// If no links are found, report and continue
		if len(scanResult.Links) == 0 {
			reporter.ReportNoLinks(file)

			continue
		}

		// Process each link
		brokenLinksFound := false
		validLinksCount := 0

		for _, link := range scanResult.Links {
			// URL-decode the link path
			decodedPath, err := url.QueryUnescape(link.Path)
			if err != nil {
				reporter.ReportProcessingError(file, err)

				continue
			}

			// Construct the full path relative to the file's location
			fullPath := fileutils.ResolveRelativePath(file, decodedPath)

			// Check if target file exists
			if !fileutils.FileExists(fullPath) {
				lineContent, _ := fileutils.GetLineContent(file, link.Line)
				reporter.ReportBrokenLink(file, link, "file not found", lineContent)

				brokenLinksFound = true

				continue
			}

			// If an anchor exists, check if it's valid
			if link.Anchor != "" {
				// Get anchors from target file (uses cache when possible)
				targetScan, err := scanner.ScanFile(fullPath)
				if err != nil {
					reporter.ReportProcessingError(file, err)

					continue
				}

				// Check if the anchor exists
				if !anchor.Contains(targetScan.Anchors, link.Anchor) {
					lineContent, _ := fileutils.GetLineContent(file, link.Line)
					reporter.ReportBrokenLink(file, link, "anchor not found", lineContent)

					brokenLinksFound = true

					continue
				}
			}

			// If we got here, the link is valid
			validLinksCount++
		}

		// Report valid links if any
		reporter.ReportValidLinks(file, validLinksCount, brokenLinksFound)

		// Update exit code if broken links were found
		if brokenLinksFound {
			exitCode = exitcode.BrokenLinks
		}
	}

	// Report success if all links are valid
	reporter.ReportSuccess()

	return exitCode
}
