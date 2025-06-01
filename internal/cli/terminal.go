package cli

import (
	"os"
)

// ColorScheme holds ANSI color codes
type ColorScheme struct {
	Bold   string
	Red    string
	Yellow string
	Green  string
	Gray   string
	Reset  string
}

const (
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	gray   = "\033[90m"
	reset  = "\033[0m"
)

// GetColorScheme returns color scheme based on terminal capabilities
func GetColorScheme(useColors bool) ColorScheme {
	if useColors {
		return ColorScheme{
			Bold:   bold,
			Red:    red,
			Yellow: yellow,
			Green:  green,
			Gray:   gray,
			Reset:  reset,
		}
	}

	return ColorScheme{"", "", "", "", "", ""}
}

// IsTerminal checks if stdout is a terminal
func IsTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
