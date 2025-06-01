package color

import (
	"os"
)

// ANSI color codes
const (
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	gray   = "\033[90m"
	reset  = "\033[0m"
)

type Palette struct {
	Bold   string
	Red    string
	Yellow string
	Green  string
	Gray   string
	Reset  string
}

func GetPalette(forceColor bool) Palette {
	useColors := isTerminal() || forceColor

	if useColors {
		return Palette{
			Bold:   bold,
			Red:    red,
			Yellow: yellow,
			Green:  green,
			Gray:   gray,
			Reset:  reset,
		}
	}

	return Palette{"", "", "", "", "", ""} // in case program is being piped into a file or another command
}

func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
