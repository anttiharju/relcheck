package cli

type Options struct {
	Verbose    bool
	ForceColor bool
}

type Command int

const (
	Usage Command = iota
	ShowVersion
	RunOnAllMarkdown
	RunOnInputFiles
)

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
