package cli

import (
	"context"

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
)

type Options struct {
	Verbose    bool
	ForceColor bool
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
	}
	inputFiles := []string{}

	for i := range args {
		arg := args[i]
		switch arg {
		case "--verbose":
			options.Verbose = true
		case "--color=always":
			options.ForceColor = true
		case "version":
			command = ShowVersion
		case "run":
			command = RunOnAllMarkdown
		default:
			command = RunOnInputFiles

			inputFiles = append(inputFiles, arg)
		}
	}

	if len(inputFiles) == 0 {
		command = Usage // fallback
	}

	return command, options, inputFiles
}
