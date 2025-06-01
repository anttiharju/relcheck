package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetLineContent retrieves a specific line from a file
func GetLineContent(filename string, lineNumber int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		if currentLine == lineNumber {
			return scanner.Text(), nil
		}
	}

	return "", fmt.Errorf("line %d not found", lineNumber)
}

// ResolveRelativePath resolves a relative path against a base file
func ResolveRelativePath(baseFile, relativePath string) string {
	dir := filepath.Dir(baseFile)
	return filepath.Join(dir, relativePath)
}
