package check

import (
	"net/url"
	"regexp"

	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/fileutils"
	"github.com/anttiharju/relcheck/internal/markdown/anchor"
	"github.com/anttiharju/relcheck/internal/markdown/link"
	"github.com/anttiharju/relcheck/internal/markdown/scan"
	"github.com/anttiharju/relcheck/internal/reporter"
)

func RelativeLinksAndAnchors(verbose, forceColors bool, files []string) exitcode.Exitcode {
	report := reporter.New(verbose, forceColors)
	exitCode := exitcode.Success

	for _, filepath := range files {
		fileExitCode := isFileValid(filepath, report)
		if fileExitCode != exitcode.Success {
			exitCode = fileExitCode
		}
	}

	if exitCode == exitcode.Success {
		report.Success()
	}

	return exitCode
}

func isFileValid(filepath string, report *reporter.Reporter) exitcode.Exitcode {
	if !fileutils.FileExists(filepath) {
		report.FileNotFound(filepath)

		return exitcode.BrokenLinks
	}

	scanResult, err := scan.File(filepath)
	if err != nil {
		report.ScanError(filepath, err)

		return exitcode.BrokenLinks
	}

	if len(scanResult.Links) == 0 {
		report.NoLinks(filepath)

		return exitcode.Success
	}

	return areLinksValid(filepath, scanResult, report)
}

func areLinksValid(filepath string, scanResult scan.Result, report *reporter.Reporter) exitcode.Exitcode {
	brokenLinksFound := false
	validLinksCount := 0

	for _, link := range scanResult.Links {
		valid := isLinkValid(filepath, link, report)
		if valid {
			validLinksCount++
		} else {
			brokenLinksFound = true
		}
	}

	report.ValidLinks(filepath, validLinksCount, brokenLinksFound)

	if brokenLinksFound {
		return exitcode.BrokenLinks
	}

	return exitcode.Success
}

func isLinkValid(filepath string, link link.Link, report *reporter.Reporter) bool {
	decodedPath, err := url.QueryUnescape(link.Path)
	if err != nil {
		report.ScanError(filepath, err)

		return false
	}

	fullpath := fileutils.ResolvePath(filepath, decodedPath)

	// If target does not exist, report it
	if !fileutils.FileExists(fullpath) {
		report.BrokenLink(filepath, link, "target not found", link.LineContent)

		return false
	}

	// If link has an anchor, validate it
	if link.Anchor != "" {
		return isAnchorValid(filepath, fullpath, link, report)
	}

	return true
}

var lineRegex = regexp.MustCompile(`^L\d+$`)

func isAnchorValid(filepath, targetpath string, link link.Link, report *reporter.Reporter) bool {
	// Get scan results for the target file once
	targetFile, err := scan.File(targetpath)
	if err != nil {
		report.ScanError(filepath, err)

		return false
	}

	// Check if it's a line link (e.g. #L5)
	if lineRegex.MatchString(link.Anchor) {
		// Validate the line number
		lineNumStr := link.Anchor[1:]

		lineNum, err := fileutils.ParseLineNumber(lineNumStr)
		if err != nil {
			report.BrokenLink(filepath, link, "invalid line number", link.LineContent)

			return false
		}

		// Check that the line number exists in the target file
		if lineNum <= 0 || lineNum > targetFile.LineCount {
			report.BrokenLink(filepath, link, "line number out of range", link.LineContent)

			return false
		}

		return true
	}

	// It's a regular anchor link
	if !anchor.Exists(targetFile.Anchors, link.Anchor) {
		report.BrokenLink(filepath, link, "heading not found", link.LineContent)

		return false
	}

	return true
}
