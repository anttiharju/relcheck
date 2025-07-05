package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/anttiharju/relcheck/internal/check"
	"github.com/anttiharju/relcheck/internal/exitcode"
	"github.com/anttiharju/relcheck/internal/git"
	"github.com/anttiharju/relcheck/internal/usage"
	"github.com/anttiharju/relcheck/internal/version"
)

type Command int

const (
	Usage Command = iota
	ShowVersion
	RunOnAllMarkdown
	RunOnInputFiles
	InvalidArgs
)

type Options struct {
	Verbose    bool
	ForceColor bool
	Directory  string
}

func Run(_ context.Context, args []string) exitcode.Exitcode {
	cmd, opts, inputFiles := ParseArgs(args)

	switch cmd {
	case Usage:
		return usage.Print()
	case ShowVersion:
		return version.Print()
	case RunOnAllMarkdown:
		return check.RelativeLinksAndAnchors(opts.Verbose, opts.ForceColor, git.ListMarkdownFiles())
	case InvalidArgs:
		return exitcode.InvalidArgs
	case RunOnInputFiles:
		fallthrough
	default:
		return check.RelativeLinksAndAnchors(opts.Verbose, opts.ForceColor, inputFiles)
	}
}

func ParseArgs(args []string) (Command, Options, []string) {
	command := RunOnInputFiles // default
	options := Options{
		Verbose:    false,
		ForceColor: false,
		Directory:  "",
	}
	inputFiles := []string{}

	index := 0
	for index < len(args) {
		arg := args[index]
		index++

		if handleOption(arg, &options, &command, &inputFiles, &index, args) {
			continue
		}
	}

	if options.Directory != "" {
		if err := os.Chdir(options.Directory); err != nil {
			fmt.Println("Error: Unable to change directory.")

			command = InvalidArgs
		}
	}

	if command == RunOnInputFiles && len(inputFiles) == 0 {
		command = Usage // fallback
	}

	return command, options, inputFiles
}

func handleOption(
	arg string,
	options *Options,
	command *Command,
	inputFiles *[]string,
	index *int,
	args []string,
) bool {
	switch arg {
	case "--verbose":
		options.Verbose = true
	case "--color=always":
		options.ForceColor = true
	case "-C", "--directory":
		if *index < len(args) {
			options.Directory = args[*index]
			*index++
		} else {
			*command = Usage
		}
	case "version", "-v", "--version":
		*command = ShowVersion
	case "all":
		*command = RunOnAllMarkdown
	default:
		*command = RunOnInputFiles

		*inputFiles = append(*inputFiles, arg)
	}

	return true
}
