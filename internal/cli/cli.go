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

type Command int

const (
	Usage Command = iota
	ShowVersion
	RunOnAllMarkdown
	RunOnInputFiles
)

type Options struct {
	Verbose    bool
	ForceColor bool
}

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
		fallthrough
	default:
		return showUsage()
	}
}

func ParseArgs(args []string) (Command, Options, []string) {
	command := Usage // default
	options := Options{
		Verbose:    false,
		ForceColor: false,
	}
	files := []string{}

	for i := range args {
		arg := args[i]
		switch arg {
		case "--verbose":
			options.Verbose = true
		case "--color=always":
			options.ForceColor = true
		case "run":
			command = RunOnAllMarkdown
		case "version":
			command = ShowVersion
		default:
			command = RunOnInputFiles

			files = append(files, arg)
		}
	}

	return command, options, files
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
