package cli

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"

	"github.com/anttiharju/relcheck/pkg/version"
)

// Options holds CLI flags and files
type Options struct {
	Verbose    bool
	ForceColor bool
	Files      []string
}

// ParseArgs parses command-line arguments to match bash script behavior
func ParseArgs(args []string) Options {
	opts := Options{
		Verbose:    false,
		ForceColor: false,
		Files:      []string{},
	}

	for i := range args {
		arg := args[i]
		switch arg {
		case "--verbose":
			opts.Verbose = true
		case "--color=always":
			opts.ForceColor = true
		case "run":
			// Use git ls-files to find all markdown files
			cmd := exec.Command("git", "ls-files", "*.md")

			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err == nil {
				scanner := bufio.NewScanner(&out)
				for scanner.Scan() {
					opts.Files = append(opts.Files, scanner.Text())
				}
			}
		case "version":
			os.Exit(version.Print("relcheck"))
		default:
			opts.Files = append(opts.Files, arg)
		}
	}

	return opts
}
