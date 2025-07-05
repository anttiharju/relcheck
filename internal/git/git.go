package git

import (
	"bytes"
	"os/exec"
	"strings"
)

func ListMarkdownFiles() []string {
	cmd := exec.Command("git", "ls-files", "-z", "*.md")
	files := []string{}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		output := out.String()
		if output != "" {
			nul := "\x00"
			for _, file := range strings.Split(output, nul) {
				if file != "" {
					files = append(files, file)
				}
			}
		}
	}

	return files
}
