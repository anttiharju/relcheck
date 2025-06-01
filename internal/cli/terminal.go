package cli

import (
	"os"

	"github.com/anttiharju/relcheck/pkg/colors"
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

// GetColorScheme returns color scheme based on terminal capabilities
func GetColorScheme(useColors bool) ColorScheme {
	if useColors {
		return ColorScheme{
			Bold:   colors.Bold,
			Red:    colors.Red,
			Yellow: colors.Yellow,
			Green:  colors.Green,
			Gray:   colors.Gray,
			Reset:  colors.Reset,
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
