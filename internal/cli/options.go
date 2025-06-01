package cli

type Options struct {
	Verbose    bool
	ForceColor bool
	Files      []string
	Command    Command
}

type Command int

const (
	CheckFiles Command = iota
	RunOnAllMarkdown
	ShowVersion
)

func ParseOptions(args []string) Options {
	opts := Options{
		Verbose:    false,
		ForceColor: false,
		Files:      []string{},
		Command:    CheckFiles,
	}

	for i := range args {
		arg := args[i]
		switch arg {
		case "--verbose":
			opts.Verbose = true
		case "--color=always":
			opts.ForceColor = true
		case "run":
			opts.Command = RunOnAllMarkdown
		case "version":
			opts.Command = ShowVersion
		default:
			opts.Files = append(opts.Files, arg)
		}
	}

	return opts
}
