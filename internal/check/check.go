package check

import (
	"net/url"

	"github.com/anttiharju/relcheck/internal/color"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/fileutils"
	"github.com/anttiharju/relcheck/internal/markdown/anchor"
	"github.com/anttiharju/relcheck/internal/markdown/scan"
	"github.com/anttiharju/relcheck/internal/reporter"
)

//nolint:cyclop,funlen
func RelativeLinksAndAnchors(verbose, forceColors bool, files []string) exitcode.Exitcode {
	// Initialize reporter with terminal colors
	colors := color.GetPalette(forceColors)
	reporter := reporter.New(colors, verbose)

	// Default exit code is success
	exitCode := exitcode.Success

	for _, filepath := range files {
		// Check if the file exists
		if !fileutils.FileExists(filepath) {
			reporter.FileNotFound(filepath)

			exitCode = exitcode.BrokenLinks

			continue
		}

		// Scan the file for links
		scanResult, err := scan.File(filepath)
		if err != nil {
			reporter.ScanError(filepath, err)

			exitCode = exitcode.BrokenLinks

			continue
		}

		// If no links are found, report and continue
		if len(scanResult.Links) == 0 {
			reporter.NoLinks(filepath)

			continue
		}

		// Process each link
		brokenLinksFound := false
		validLinksCount := 0

		for _, link := range scanResult.Links {
			// URL-decode the link path
			decodedPath, err := url.QueryUnescape(link.Path)
			if err != nil {
				reporter.ScanError(filepath, err)

				continue
			}

			// Construct the full path relative to the file's location
			fullpath := fileutils.ResolveRelativePath(filepath, decodedPath)

			// Check if target file exists
			if !fileutils.FileExists(fullpath) {
				lineContent, _ := fileutils.GetLineContent(filepath, link.Line)
				reporter.BrokenLink(filepath, link, "file not found", lineContent)

				brokenLinksFound = true

				continue
			}

			// If an anchor exists, check if it's valid
			if link.Anchor != "" {
				// Get anchors from target file (uses cache when possible)
				targetScan, err := scan.File(fullpath)
				if err != nil {
					reporter.ScanError(filepath, err)

					continue
				}

				// Check if the anchor exists
				if !anchor.Contains(targetScan.Anchors, link.Anchor) {
					lineContent, _ := fileutils.GetLineContent(filepath, link.Line)
					reporter.BrokenLink(filepath, link, "anchor not found", lineContent)

					brokenLinksFound = true

					continue
				}
			}

			// If we got here, the link is valid
			validLinksCount++
		}

		// Report valid links if any
		reporter.ValidLinks(filepath, validLinksCount, brokenLinksFound)

		// Update exit code if broken links were found
		if brokenLinksFound {
			exitCode = exitcode.BrokenLinks
		}
	}

	// Report success if all links are valid
	reporter.Success()

	return exitCode
}
