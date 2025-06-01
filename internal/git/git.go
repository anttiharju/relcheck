package git

import (
	"bufio"
	"bytes"
	"os/exec"
)

func ListMarkdownFiles() []string {
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
