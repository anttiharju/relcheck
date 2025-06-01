package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

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

func ResolvePath(baseFile, relativePath string) string {
	dir := filepath.Dir(baseFile)

	return filepath.Join(dir, relativePath)
}
