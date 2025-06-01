package fileutils

import (
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func ResolvePath(baseFile, relativePath string) string {
	dir := filepath.Dir(baseFile)

	return filepath.Join(dir, relativePath)
}
