package cli

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/program"
	"github.com/anttiharju/relcheck/internal/version"
)

// ExecuteCommand handles special commands like "run" and "version"
func Run(_ context.Context, args []string) exitcode.Exitcode {
	cmd, opts, files := ParseArgs(args)
	switch cmd {
	case ShowVersion:
		return version.Print("relcheck")
	case RunOnAllMarkdown:
		return program.Start(opts.Verbose, opts.ForceColor, findAllMarkdownFiles())
	case RunOnInputFiles:
		return program.Start(opts.Verbose, opts.ForceColor, files)
	case Usage:
		return showUsage()
	default:
		return showUsage()
	}
}

// findAllMarkdownFiles uses git ls-files to find all markdown files
func findAllMarkdownFiles() []string {
	cmd := exec.Command("git", "ls-files", "*.md")
	files := []string{}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
	}

	return files
}

func showUsage() exitcode.Exitcode {
	fmt.Println("Usage: relcheck [--verbose] [--color=always] <file1.md> [file2.md] ...")
	fmt.Println("   or: relcheck [--verbose] [--color=always] run  (to check all *.md files in Git)")
	fmt.Println("   or: relcheck version  (to show version information)")

	return exitcode.UsageError
}
