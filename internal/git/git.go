package git

import (
	"bytes"
	"context"
	"os/exec"
)

func ListMarkdownFiles(ctx context.Context) []string {
	out, err := exec.CommandContext(ctx, "git", "ls-files", "-z", "*.md").Output()
	if err != nil {
		return nil
	}

	// For empty output
	if len(out) == 0 {
		return nil
	}

	// Split by NUL byte - the output has a trailing NUL that creates an empty final element
	parts := bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0})

	// Convert all parts to strings
	files := make([]string, len(parts))
	for i, part := range parts {
		files[i] = string(part)
	}

	return files
}
