package cli

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"

	"github.com/anttiharju/relcheck/pkg/version"
)

// ExecuteCommand handles special commands like "run" and "version"
func ExecuteCommand(opts *Options) bool {
	switch opts.Command {
	case RunOnAllMarkdown:
		findAllMarkdownFiles(opts)

		return true
	case ShowVersion:
		showVersion()

		return false
	case CheckFiles:
		return true
	default:
		return true
	}
}

// findAllMarkdownFiles uses git ls-files to find all markdown files
func findAllMarkdownFiles(opts *Options) {
	cmd := exec.Command("git", "ls-files", "*.md")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			opts.Files = append(opts.Files, scanner.Text())
		}
	}
}

// showVersion displays version information and exits
func showVersion() {
	os.Exit(version.Print("relcheck"))
}
