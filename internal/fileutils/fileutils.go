package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func ResolvePath(baseFile, relativePath string) string {
	dir := filepath.Dir(baseFile)

	return filepath.Join(dir, relativePath)
}

func ParseLineNumber(s string) (int, error) {
	lineNum, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %s", s)
	}

	return lineNum, nil
}

func CountLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error counting lines: %w", err)
	}

	return lineCount, nil
}
